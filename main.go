package main

import "fmt"

type Task struct {
	ID          int
	Description string
	Status      bool
}

var tasks []Task

func main() {
	addTask("Template one")
	addTask("Template two")
	addTask("Template three")
	viewTasks()
}

func addTask(description string) {
	id := len(tasks) + 1
	task := Task{
		ID:          id,
		Description: description,
		Status:      false,
	}
	tasks = append(tasks, task)
	fmt.Println("Task added successfully")
}

func viewTasks() {
	fmt.Print("ID\t|Status\t|\tDescription\n")
	for _, task := range tasks {
		if len(tasks) != 0 {
			fmt.Printf("%v\t|%v\t|%v\n", task.ID, task.Status, task.Description)
		} else {
			fmt.Println("No tasks for now")
		}
	}
}
