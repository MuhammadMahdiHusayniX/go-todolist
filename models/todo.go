package models

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Todo struct {
	gorm.Model
	Task      string
	CreatedBy string
	Completed bool
}

func Setup() {
	var err error
	db, err = gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("models setup error: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&Todo{})
}

func GetTodos() ([]Todo, error) {
	var (
		todos []Todo
		err   error
	)

	err = db.Find(&todos).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return todos, nil
}

func AddTodo(task string, createdBy string) (uint, error) {
	todo := Todo{
		Task:      task,
		CreatedBy: createdBy,
		Completed: false,
	}

	if err := db.Create(&todo).Error; err != nil {
		return 0, err
	}

	return todo.ID, nil
}

func DeleteTodo(id int) error {
	if err := db.Where("id = ?", id).Delete(&Todo{}).Error; err != nil {
		return err
	}

	return nil
}

func MarkComplete(id int) error {
	if err := db.Model(&Todo{}).Where("id = ?", id).Update("Completed", true).Error; err != nil {
		return err
	}

	return nil
}
