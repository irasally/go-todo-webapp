package main

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"github.com/gotailwindcss/tailwind/twembed"
	"github.com/gotailwindcss/tailwind/twhandler"
	"github.com/labstack/echo/v4"

	"todo-webapp/db"
	"todo-webapp/model"
	"todo-webapp/router"
)

//go:embed static
var static embed.FS

//go:embed templates
var templates embed.FS

type Tempalte struct {
	templates *template.Template
}

func (t *Tempalte) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func formatDateTime(d time.Time) string {
	if d.IsZero() {
		return ""
	}
	return d.Format("2006-01-02 15:04")
}

func main() {
	db.CreateDBConnection()
	model.SetupTodoTable()
	defer db.CloseDBConnection()

	e := echo.New()
	router.SetRouter(e)

	e.Renderer = &Tempalte{
		templates: template.Must(template.New("").
		Funcs(template.FuncMap{
			"FormatDateTime": formatDateTime,
		}).ParseFS(templates, "templates/*")),
	}

	staticFs, err := fs.Sub(static, "static")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FileSystem(http.FS(staticFs)))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", fileServer)))

	e.GET("/css/*", echo.WrapHandler(
										http.StripPrefix("/css/", twhandler.New(
																									http.Dir("css"),
																									"/css",
																									twembed.New()))))
	e.Logger.Fatal(e.Start(":8989"))
}
