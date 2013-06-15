package controllers

import (
	"fmt"
	"github.com/robfig/revel"
	"github.com/robfig/revel/cache"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Check struct {
	Passed  bool
	Comment string
}

type Review struct {
	Checks []Check
}

func (r *Review) ChecksPassed() int {
	count := 0
	for _, check := range r.Checks {
		if check.Passed {
			count = count + 1
		}

	}
	return count
}

type Repository struct {
	*Review
	Description string
	Name        string
	User        string
}

type App struct {
	*revel.Controller
}

func FetchRepository(user string, repo string) Repository {
	var review Review
	key := "github-" + user + "-" + repo
	err := cache.Get(key, &review)

	r := Repository{Description: "", Name: repo, User: user}

	if err != nil {
		return r
	}

	r.Review = &review
	return r
}

func FindFiles(user string, repo string, files ...string) bool {
	base := "https://api.github.com"

	found := make(chan bool)
	misses := make(chan bool)

	for _, file := range files {
		go func(b string, u string, r string, f string) {
			url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", b, u, r, f)
			resp, err := http.Get(url)

			if err != nil {
				misses <- true
				return
			}

			if resp.StatusCode == 200 {
				found <- true
				return
			}

			misses <- true
		}(base, user, repo, file)
	}

	count := 0

	select {
	case <-found:
		return true
	case <-misses:
		count = count + 1
		if count == len(files) {
			return false
		}
	}

	return false
}

func (c App) Index() revel.Result {
	repos := []Repository{
		FetchRepository("mooseburger", "AbcaTSH"),
		FetchRepository("blakeembrey", "code-problems"),
		FetchRepository("m242", "maildrop"),
		FetchRepository("qq99", "echoplexus"),
		FetchRepository("lampepfl", "scala-js"),
		FetchRepository("creaktive", "rainbarf"),
		FetchRepository("AFNetworking", "AFNetworking"),
		FetchRepository("mitsuhiko", "flask"),
		FetchRepository("libgit2", "libgit2"),
	}

	return c.Render(repos)
}

func (c App) FindReport(repository string) revel.Result {
	c.Validation.Required(repository)
	c.Validation.Match(repository, regexp.MustCompile("^[^/]+/[^/]+$"))

	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Index)
	}

	lines := strings.Split(repository, "/")
	return c.Redirect("/github/%s/%s", lines[0], lines[1])
}

func (c App) ReportCard(user string, repo string) revel.Result {
	title := user + "/" + repo + "'s Report Card - Rotten Repositories"
	key := "github-" + user + "-" + repo
	partial := c.Request.Header.Get("X-PJAX") != ""

	// Create a review channel
	var review Review

	err := cache.Get(key, &review)
	missing := err != nil

	if missing && partial {

		checks := make(chan Check)
		base := "https://api.github.com"

		go func() {
			repo_url := fmt.Sprintf("%s/repos/%s/%s", base, user, repo)
			resp, err := http.Get(repo_url)

			if err != nil {
				checks <- Check{Passed: false, Comment: "GitHub Error :("}
				return
			}

			if resp.StatusCode != 200 {
				checks <- Check{Passed: false, Comment: "Repository doesn't exist"}
				return
			}

			checks <- Check{Passed: true, Comment: "Repository exists"}
			return
		}()

		// Check for LICENSE
		go func() {
			if FindFiles(user, repo, "LICENSE", "LICENSE.txt") {
				checks <- Check{Passed: true, Comment: "LICENSE in repository"}
			} else {
				checks <- Check{Passed: false, Comment: "No LICENSE in repository"}
			}
		}()

		// Check for README
		go func() {
			readme_url := fmt.Sprintf("%s/repos/%s/%s/readme", base, user, repo)
			resp, err := http.Get(readme_url)

			if err != nil {
				checks <- Check{Passed: false, Comment: "GitHub error :("}
				return
			}

			if resp.StatusCode != 200 {
				checks <- Check{Passed: false, Comment: "No README in repository"}
				return
			}

			checks <- Check{Passed: true, Comment: "README in repository"}
		}()

		for i := 0; i < 3; i++ {
			review.Checks = append(review.Checks, <-checks)
		}

		go cache.Set(key, review, 24*time.Hour)
	}

	return c.Render(user, repo, title, partial, review, missing)
}
