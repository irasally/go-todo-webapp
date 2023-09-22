package router

import (
	"errors"
	"net/http"
	"todo-webapp/model"

	"github.com/labstack/echo/v4"
)

type Data struct {
	Todos []model.Todo
	Errors []error
}

func root_get(e *echo.Echo){
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
}

func root_post(e *echo.Echo) {
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
			err = add(todo)
		} else {
			if c.FormValue("delete") != "" {
				err = delete(todo)
			} else {
				err = update(todo)
			}
		}
		if err != nil {
			return c.Render(http.StatusBadRequest, "index", Data{Errors: []error{err}})
		}
		return c.Redirect(http.StatusFound, "/")
	})
}

func add(todo model.Todo) error {
	err := model.AddTodo(todo)
	if err != nil {
		logger.Error(err)
		err = errors.New("Cannot insert todo")
	}
	return err
}

func delete(todo model.Todo) error{
	err := model.DeleteTodo(todo)
	if err != nil {
		logger.Error(err)
		err = errors.New("Cannot delete todo")
	}
	return err
}

func update(todo model.Todo) error {
	err := model.UpdateTodo(todo)
	if err != nil {
		logger.Error(err)
		err = errors.New("Cannot update todo")
	}
	return err
}
