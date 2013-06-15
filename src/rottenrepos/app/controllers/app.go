package controllers

import "github.com/robfig/revel"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) ReportCard(user string, repo string) revel.Result {
	title := user + "/" + repo + "'s Report Card - Rotten Repositories"
	partial := c.Request.Header.Get("X-PJAX") != ""

	if partial {
		//lots of work
	}

	return c.Render(user, repo, title, partial)
}
