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

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"

	"todo-webapp/model"
	"todo-webapp/router"
)

//go:embed static
var static embed.FS

//go:embed templates
var templates embed.FS

type Todo struct {
	bun.BaseModel `bun:"table:todos,alias:t"`

	ID int64 `bun:"id,pk,autoincrement"`
	Content string `bun:"content,notnull"`
	Done bool `bun:"done"`
	Until time.Time `bun:"until,nullzero"`
	CreatedAt time.Time
	UpdatedAt time.Time `bun:",nullzero"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type Data struct {
	Todos []Todo
	Errors []error
}

func customFunc(todo *Todo) func([]string) []error{
	return func(values []string) []error {
		if len(values) == 0 || values[0] == "" {
			return nil
		}
		dt, err := time.Parse("2006-01-02T15:04 MST", values[0]+" JST")
		if err != nil {
			return []error{err}
		}
		todo.Until = dt
		return nil
	}
}

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
	db := model.DBConnection()
	model.CreateTable(db)

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
	e.Logger.Fatal(e.Start(":8989"))
}
