package router

import (
	"context"
	"errors"
	"net/http"
)

func SetRouter(e *echo.Echo) error {
	e.GET("/", func(c echo.Context) error {
		var todos []Todo
		ctx := context.Background()
		err := db.NewSelect().Model(&todos).Order("created_at").Scan(ctx)
		if err != nil {
			e.Logger.Error(err)
			return c.Render(http.StatusBadRequest, "index", Data{
				Errors: []error{errors.New("Cannot get todos")},
			})
		}
		return c.Render(http.StatusOK, "index", Data{ Todos: todos })
	})
	e.POST("/", func(c echo.Context) error {
		var todo Todo
		// パラメーターをフィールドにバインド
		errs := echo.FormFieldBinder(c).
		Int64("id", &todo.ID).
		String("content", &todo.Content).
		Bool("done", &todo.Done).
		CustomFunc("until", customFunc(&todo)).
		BindErrors()
		if errs != nil {
			e.Logger.Error(err)
			return c.Render(http.StatusBadRequest, "index", Data{Errors: errs})
		} else if todo.ID == 0 {
			// IDが0の特は登録
			ctx := context.Background()
			if todo.Content == "" {
				err = errors.New("Content is empty")
			} else {
				_, err = db.NewInsert().Model(&todo).Exec(ctx)
				if err != nil {
					e.Logger.Error(err)
					err = errors.New("Cannot insert todo")
				}
			}
		} else {
			ctx := context.Background()
			if c.FormValue("delete") != "" {
				// 削除
				_, err = db.NewDelete().Model(&todo).Where("id = ?", todo.ID).Exec(ctx)
			} else {
				// 更新
				var orig Todo
				err = db.NewSelect().Model(&orig).Where("id = ?", todo.ID).Scan(ctx)
				if err == nil {
					orig.Done = todo.Done
					_, err = db.NewUpdate().Model(&orig).Where("id = ?", todo.ID).Exec(ctx)
				}
			}
			if err != nil {
				e.Logger.Error(err)
				err = errors.New("Cannot update todo")
			}
		}
		if err != nil {
			return c.Render(http.StatusBadRequest, "index", Data{Errors: []error{err}})
		}
		return c.Redirect(http.StatusFound, "/")
	})

}
