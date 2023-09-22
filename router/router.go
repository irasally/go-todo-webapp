package router

import (
	"errors"
	"net/http"
	"time"

	"todo-webapp/model"

	"github.com/labstack/echo/v4"
)

type Data struct {
	Todos []model.Todo
	Errors []error
}

func customFunc(todo *model.Todo) func([]string) []error{
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

func setupLogger(e *echo.Echo){
	logger = e.Logger
}

var logger echo.Logger

func SetRouter(e *echo.Echo) {
	setupLogger(e)

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
				err = doneTodo(todo)
			}
		}

		if err != nil {
			return c.Render(http.StatusBadRequest, "index", Data{Errors: []error{err}})
		}
		return c.Redirect(http.StatusFound, "/")
	})
}

func doneTodo(todo model.Todo) error {
	// 更新
	err := model.UpdateTodo(todo)
	if err != nil {
		logger.Error(err)
		err = errors.New("Cannot update todo")
	}
	return err
}
