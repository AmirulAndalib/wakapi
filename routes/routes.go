package routes

import (
	"fmt"
	"github.com/muety/wakapi/config"
	"github.com/muety/wakapi/models"
	"github.com/muety/wakapi/utils"
	"html/template"
	"io/ioutil"
	"path"
	"strings"
)

func Init() {
	loadTemplates()
}

var templates map[string]*template.Template

func loadTemplates() {
	tplPath := "views"
	tpls := template.New("").Funcs(template.FuncMap{
		"json":        utils.Json,
		"date":        utils.FormatDateHuman,
		"title":       strings.Title,
		"join":        strings.Join,
		"add":         utils.Add,
		"capitalize":  utils.Capitalize,
		"toRunes":     utils.ToRunes,
		"entityTypes": models.SummaryTypes,
		"typeName":    typeName,
		"getBasePath": func() string {
			return config.Get().Server.BasePath
		},
		"getVersion": func() string {
			return config.Get().Version
		},
		"getDbType": func() string {
			return strings.ToLower(config.Get().Db.Type)
		},
		"htmlSafe": func(html string) template.HTML {
			return template.HTML(html)
		},
	})
	templates = make(map[string]*template.Template)

	files, err := ioutil.ReadDir(tplPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		tplName := file.Name()
		if file.IsDir() || path.Ext(tplName) != ".html" {
			continue
		}

		tpl, err := tpls.New(tplName).ParseFiles(fmt.Sprintf("%s/%s", tplPath, tplName))
		if err != nil {
			panic(err)
		}

		templates[tplName] = tpl
	}
}

func typeName(t uint8) string {
	if t == models.SummaryProject {
		return "project"
	}
	if t == models.SummaryLanguage {
		return "language"
	}
	if t == models.SummaryEditor {
		return "editor"
	}
	if t == models.SummaryOS {
		return "operating system"
	}
	if t == models.SummaryMachine {
		return "machine"
	}
	return "unknown"
}
