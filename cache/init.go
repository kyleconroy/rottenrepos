package cache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
)

var instance redis.Conn

// Serialize transforms the given value into bytes following these rules:
// - If value is a byte array, it is returned as-is.
// - If value is an int or uint type, it is returned as the ASCII representation
// - Else, encoding/gob is used to serialize
func Serialize(value interface{}) ([]byte, error) {
	if bytes, ok := value.([]byte); ok {
		return bytes, nil
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(v.Uint(), 10)), nil
	}

	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	if err := encoder.Encode(value); err != nil {
		log.Printf("revel/cache: gob encoding '%s' failed: %s", value, err)
		return nil, err
	}
	return b.Bytes(), nil
}

// Deserialize transforms bytes produced by Serialize back into a Go object,
// storing it into "ptr", which must be a pointer to the value type.
func Deserialize(byt []byte, ptr interface{}) (err error) {
	if bytes, ok := ptr.(*[]byte); ok {
		*bytes = byt
		return
	}

	if v := reflect.ValueOf(ptr); v.Kind() == reflect.Ptr {
		switch p := v.Elem(); p.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var i int64
			i, err = strconv.ParseInt(string(byt), 10, 64)
			if err != nil {
				log.Printf("revel/cache: failed to parse int '%s': %s", string(byt), err)
			} else {
				p.SetInt(i)
			}
			return

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var i uint64
			i, err = strconv.ParseUint(string(byt), 10, 64)
			if err != nil {
				log.Printf("revel/cache: failed to parse uint '%s': %s", string(byt), err)
			} else {
				p.SetUint(i)
			}
			return
		}
	}

	b := bytes.NewBuffer(byt)
	decoder := gob.NewDecoder(b)
	if err = decoder.Decode(ptr); err != nil {
		log.Printf("revel/cache: gob decoding failed: %s", err)
		return
	}
	return
}

func Get(key string, ptrValue interface{}) (error, bool) {
	if instance == nil {
		return errors.New("Global client hasn't been initialized"), false
	}

	val, err := instance.Do("GET", key)

	if err != nil {
		return err, false
	}

	if val == nil {
		return nil, false
	}

	bulk := val.([]byte)

	err = Deserialize(bulk, ptrValue)

	if err != nil {
		return err, false
	}

	return nil, true
}

func Set(key string, value interface{}, expires time.Duration) error {
	if instance == nil {
		return errors.New("Global client hasn't been initialized")
	}

	bulk, err := Serialize(value)

	if err != nil {
		return err
	}

	val, err := instance.Do("SETEX", key, expires.Seconds(), bulk)

	log.Println(val.(string))

	if err != nil {
		return err
	}

	return nil
}

func Init() {

	addr := "redis://localhost:6379"

	if os.Getenv("REDISCLOUD_URL") != "" {
		addr = os.Getenv("REDISCLOUD_URL")
	}

	c, err := redisurl.ConnectToURL(addr)

	if err != nil {
		log.Fatal(err)
	}

	instance = c
}
