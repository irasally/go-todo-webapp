package router

import (
	"context"
	"errors"
	"net/http"
)

func SetRouter(e *echo.Echo) error {
	e.GET("/", func(c echo.Context) error {
		todos, err := model.GetTodos()

		if err != nil {
			e.Logger.Error(err)
			return c.Render(http.StatusBadRequest, "index", Data{
				Errors: []error{errors.New("Cannot get todos")},
			})
		}
		return c.Render(http.StatusOK, "index", Data{ Todos: todos })
	})

	e.POST("/", func(c echo.Context) error {
		var todo model.Todo
		var err error

		// パラメーターをフィールドにバインド
		errs := echo.FormFieldBinder(c).
		Int64("id", &todo.ID).
		String("content", &todo.Content).
		Bool("done", &todo.Done).
		CustomFunc("until", customFunc(&todo)).
		BindErrors()

		if errs != nil {
			e.Logger.Error(errs)
			return c.Render(http.StatusBadRequest, "index", Data{Errors: errs})
		} else if todo.ID == 0 {
			err := model.AddTodo(todo)
			if err != nil {
				e.Logger.Error(err)
				err = errors.New("Cannot insert todo")
			}
		} else {
			if c.FormValue("delete") != "" {
				// 削除
				err := model.DeleteTodo(todo)
				if err != nil {
					e.Logger.Error(err)
					err = errors.New("Cannot delete todo")
				}
			} else {
				// 更新
				err := model.UpdateTodo(todo)
				if err != nil {
					e.Logger.Error(err)
					err = errors.New("Cannot update todo")
				}
			}
		}

		if err != nil {
			return c.Render(http.StatusBadRequest, "index", Data{Errors: []error{err}})
		}
		return c.Redirect(http.StatusFound, "/")
	})

}
