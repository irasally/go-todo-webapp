package model

import (
	"context"
	"errors"
	"log"
	"time"

	"todo-webapp/db"

	"github.com/uptrace/bun"
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

func GetTodos() ([]Todo, error) {
	var todos []Todo
	ctx := context.Background()
	err := db.Connection.NewSelect().Model(&todos).Order("done", "until", "created_at").Scan(ctx)

	return todos, err
}

func AddTodo(todo Todo) error {
	var err error

	// IDが0の特は登録
	ctx := context.Background()
	if todo.Content == "" {
		err = errors.New("Content is empty")
	} else {
		_, err = db.Connection.NewInsert().Model(&todo).Exec(ctx)
	}
	return err
}

func DeleteTodo(todo Todo) error {
	ctx := context.Background()
	// 削除
	_, err := db.Connection.NewDelete().Model(&todo).Where("id = ?", todo.ID).Exec(ctx)

	return err
}

func UpdateTodo(todo Todo) error {
	ctx := context.Background()

	// 更新
	var orig Todo
	err := db.Connection.NewSelect().Model(&orig).Where("id = ?", todo.ID).Scan(ctx)
	if err == nil {
		orig.Done = todo.Done
		_, err = db.Connection.NewUpdate().Model(&orig).Where("id = ?", todo.ID).Exec(ctx)
	}

	return err
}

// Task型のテーブルを作成する
func SetupTodoTable() {
	var err error
	ctx := context.Background()
	_, err = db.Connection.NewCreateTable().Model((*Todo)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
