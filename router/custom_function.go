package router

import (
	"time"
	"todo-webapp/model"
)

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
