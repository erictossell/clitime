package db

import (
	"database/sql"
	"log"
	"time"

	// Import the SQLite3 driver
	_ "modernc.org/sqlite"
	// Import the SQLite3 driver
)

type Task struct {
	ID          uint `sql:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
	Name        string
	Description string
	EndTime     time.Time
	ResetTime   time.Time
	ElapsedTime time.Duration
}

var db *sql.DB

// Changed: Initialize DB connection and table creation using raw SQL
func InitDB() {
	dsn := "./task.db"
	var err error
	db, err = sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		CreatedAt DATETIME,
		UpdatedAt DATETIME,
		DeletedAt, DATETIME,
		Name TEXT,
		Description TEXT,
		EndTime DATETIME,
		ResetTime DATETIME,
		ElapsedTime INTEGER
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

// Changed: Create a new task using raw SQL
func CreateTask(name, description string) *Task {
	insertQuery := `INSERT INTO tasks (CreatedAt, UpdatedAt, Name, Description) VALUES (datetime('now'), datetime('now'), ?, ?)`
	_, err := db.Exec(insertQuery, name, description)
	if err != nil {
		log.Fatalf("Failed to insert a new task: %v", err)
		return nil
	}

	return &Task{Name: name, Description: description}
}

// Changed: Get all tasks using raw SQL
func GetAllTasks() []Task {
	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		log.Fatalf("Failed to get all tasks: %v", err)
		return nil
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt, &task.Name, &task.Description, &task.EndTime, &task.ResetTime, &task.ElapsedTime)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
			return nil
		}
		tasks = append(tasks, task)
	}

	return tasks
}

// Changed: Get task by ID using raw SQL
func GetTaskByID(id uint) *Task {
	row := db.QueryRow("SELECT * FROM tasks WHERE ID = ?", id)
	var task Task
	err := row.Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt, &task.Name, &task.Description, &task.EndTime, &task.ResetTime, &task.ElapsedTime)
	if err != nil {
		log.Fatalf("Failed to get task by ID: %v", err)
		return nil
	}

	return &task
}

// Changed: Update task using raw SQL
func UpdateTask(id uint, newName string, newDescription string, elapsedTime time.Duration) {
	updateQuery := "UPDATE tasks SET Name = ?, Description = ?, ElapsedTime = ?, UpdatedAt = datetime('now') WHERE ID = ?"
	_, err := db.Exec(updateQuery, newName, newDescription, elapsedTime, id)
	if err != nil {
		log.Fatalf("Failed to update task: %v", err)
	}
}

// Changed: Update task elapsed time using raw SQL
func UpdateTaskElapsedTime(id uint, elapsedTime time.Duration) {
	updateQuery := "UPDATE tasks SET ElapsedTime = ?, UpdatedAt = datetime('now') WHERE ID = ?"
	_, err := db.Exec(updateQuery, elapsedTime, id)
	if err != nil {
		log.Fatalf("Failed to update task elapsed time: %v", err)
	}
}

// Changed: Get latest task using raw SQL
func GetLatestTask() *Task {
	row := db.QueryRow("SELECT * FROM tasks ORDER BY CreatedAt DESC LIMIT 1")
	var task Task
	err := row.Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt, &task.DeletedAt, &task.Name, &task.Description, &task.EndTime, &task.ResetTime, &task.ElapsedTime)
	if err != nil {

		if err == sql.ErrNoRows {
			return CreateTask("New Task", "")
		}
		log.Fatalf("Failed to get the latest task: %v", err)
		return nil
	}

	return &task
}
