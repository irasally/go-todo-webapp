package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

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

func main() {
	e := echo.New()
	e.get("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})
	e.Logger.Fatal(e.Start(":8989"))
}
