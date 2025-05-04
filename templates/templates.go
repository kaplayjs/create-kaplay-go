package templates

import (
	"github.com/fatih/color"
)

type Template struct {
	Url         string
	Name        string
	Description string
}

var DefaultTemplates = map[string]Template{
	"vanilla": {
		Url:         "https://github.com/lajbel/template-html.git",
		Name:        color.HiYellowString("vanilla"),
		Description: "HTML and JS using script import",
	},
	"vite": {
		Url:         "https://github.com/lajbel/template-vite.git",
		Name:        color.HiYellowString("vite"),
		Description: "Vite + JS",
	},
	"vite-ts": {
		Url:         "https://github.com/kaplayjs/template-vite-ts.git",
		Name:        color.HiBlueString("vite-ts"),
		Description: "Vite + TS",
	},
}
