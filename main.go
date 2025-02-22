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
	taskComplete(3)
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

func taskComplete(id int) {
	found := false
	if len(tasks) == 0 {
		fmt.Println("\tNo tasks here yet")
	} else {
		for i, task := range tasks {
			if task.ID == id {
				tasks[i].Status = true
				fmt.Printf("Task %v Marked completed\n", task.ID)
				found = true
				break
			}
		}
		if !found {
			fmt.Println("\tTask not found")
		}
	}
}

func viewTasks() {
	fmt.Print("ID\t|Status\t|\tDescription\n")
	if len(tasks) > 0 {
		for _, task := range tasks {
			fmt.Printf("%v\t|%v\t|%v\n", task.ID, task.Status, task.Description)
		}
	} else {
		fmt.Println("\tNo tasks for now")
	}

}
