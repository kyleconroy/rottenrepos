package controllers

import (
	"fmt"
	"github.com/robfig/revel"
	"github.com/robfig/revel/cache"
	"net/http"
	"time"
	"strings"
	"regexp"
)

type Check struct {
	Passed  bool
	Comment string
}

type Review struct {
	Checks []Check
}

type App struct {
	*revel.Controller
}


func (c App) Index() revel.Result {
	return c.Render()
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

		for i := 0; i < 2; i++ {
			review.Checks = append(review.Checks, <-checks)
		}

		go cache.Set(key, review, 24*time.Hour)
	}

	return c.Render(user, repo, title, partial, review, missing)
}
