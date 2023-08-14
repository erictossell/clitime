package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
  "time"
)

type Task struct {
	gorm.Model
	Name        string
	Description string
  EndTime time.Time
  ResetTime time.Time
  ElapsedTime time.Duration
}

var db *gorm.DB

func InitDB() {
	dsn := "./task.db"
	var err error
	db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&Task{})
}

func CreateTask(name, description string) *Task {
  task := &Task{Name: name, Description: description}
	db.Create(task)
	return task
}

func GetAllTasks() []Task {
	var tasks []Task
	db.Find(&tasks)
	return tasks
}

func GetTaskByID(id uint) *Task {
	var task Task
	db.First(&task, id)
	return &task
}

func UpdateTask(id uint, newName string, newDescription string, elapsedTime time.Duration) {
	var task Task
	db.First(&task, id)
  task.Name = newName
  task.Description = newDescription
  task.ElapsedTime = elapsedTime
  db.Save(&task)
}

func UpdateTaskElapsedTime(id uint, elapsedTime time.Duration) {
  var task Task
  db.First(&task, id)
  task.ElapsedTime = elapsedTime
  db.Save(&task)
}

func GetLatestTask() *Task {
  var task Task
  db.Order("created_at desc").First(&task)
  return &task
}
