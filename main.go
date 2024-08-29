package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const taskFile = "task.json"

type Task struct {
	ID          int
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func loadTasks() ([]Task, error) {
	if _, err := os.Stat(taskFile); os.IsNotExist(err) {
		return []Task{}, nil
	}

	data, err := os.ReadFile(taskFile)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(taskFile, data, 0644)
}

func addTask(descriptionParts ...string) error {
	description := strings.Join(descriptionParts, " ")
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	newTask := Task{
		ID:          len(tasks) + 1,
		Description: description,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tasks = append(tasks, newTask)

	return saveTasks(tasks)
}

func listTask(filter string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if filter == "all" || filter == task.Status {
			fmt.Printf("ID: %v\nОписание: %v\nСтатус: %v\nДата создания: %v\nДата обновления: %v\n\n",
				task.ID, task.Description, task.Status, task.CreatedAt.Format(time.RFC1123), task.UpdatedAt.Format(time.RFC1123))
		}
	}
	return nil
}

func updateTask(id int, newDescription string, newStatus string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Description = newDescription
			tasks[i].Status = newStatus
			tasks[i].UpdatedAt = time.Now()
			saveTasks(tasks)
			return nil
		}
	}
	return fmt.Errorf("задача с id %d не найдена", id)
}

func deleteTask(id int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			err = saveTasks(tasks)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("задача с id %d не найдена", id)
}

func changeId(oldId int, newId int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID == oldId {
			tasks[i].ID = newId
			saveTasks(tasks)
			return nil
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Ожидается одна из команд: 'add', 'list', 'update' или 'delete'")
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Ожидается описание задачи")
			return
		}
		description := os.Args[2:]
		err := addTask(description...)
		if err != nil {
			fmt.Println("Ошибка при добавлении задачи", err)
		} else {
			fmt.Println("Задача добавлена")
		}
	case "list":
		filter := "all"
		if len(os.Args) >= 3 {
			filter = os.Args[2]
		}
		err := listTask(filter)
		if err != nil {
			fmt.Println("Ошибка при выводе списка задач", err)
		}
	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Ожидается новый ID, описание и статус задачи")
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Некоректный ID", err)
			return
		}
		newDescription := os.Args[3]
		newStatus := os.Args[4]

		err = updateTask(id, newDescription, newStatus)
		if err != nil {
			fmt.Println("Ошибка при обновлении задачи", err)
		} else {
			fmt.Println("Задача успешно обновлена!")
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Ожидается ID задачи")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Некоректный ID", err)
			return
		}

		err = deleteTask(id)
		if err != nil {
			fmt.Println("Ошибка при удалении задачи", err)
		} else {
			fmt.Println("Задача успешно удалена")
		}
	case "change":
		if len(os.Args) < 4 {
			fmt.Println("Ожидается ввод старого id и нового id")
			return
		}
		oldId, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Неверный старый id:", os.Args[2])
		}

		newId, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("Неверный новый id", os.Args[3])
		}

		err = changeId(oldId, newId)
		if err != nil {
			fmt.Println("Ошибка смены id", err)
		} else {
			fmt.Println("Смена id прошла успешно!")
		}

	default:
		fmt.Println("Неизвестная команда. Ожидается одна из команд: 'add', 'list', 'update' или 'delete'")
	}
}
