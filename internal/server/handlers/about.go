package handlers

import (
	"net/http"

	"github.com/throskam/memo/internal/views/pages"
)

type AboutController struct{}

func NewAboutController() *AboutController {
	return &AboutController{}
}

func (c *AboutController) PageGet(w http.ResponseWriter, r *http.Request) {
	Render(w, r, pages.AboutPage())
}