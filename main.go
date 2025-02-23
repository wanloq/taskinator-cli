package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

var tasks []Task

func main() {
	/*
	 ADD-Go to church on sunday
	 ADD-Read a book
	 ADD-Read the bible
	 ADD-Drink water
	 ADD-Go to the Gym at 15:00
	*/
	loadTasks()
	commands := map[string]string{
		"Add new task":           "ADD-\"Task name\"",
		"Delete a task":          "DELETE-\"Task ID\"",
		"View all tasks":         "VIEW ALL",
		"View Completed tasks":   "VIEW COMPLETED",
		"View Pending tasks":     "VIEW PENDING",
		"Mark task as Completed": "MARK COMPLETE-\"Task ID\"",
	}

	for {
		fmt.Println("\n\tCLI TASKINATOR IS RUNNING")
		fmt.Print("\nTo perform a command, Please enter the prompt below without the quotation marks (e.g. ADD-Read a book) \n\n")
		for command, prompt := range commands {
			fmt.Printf("To %v:: \t%v\n", command, prompt)
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nEnter a command: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		inputSlice := strings.Split(input, "-")
		command := strings.ToUpper(inputSlice[0])
		var arg string
		if len(inputSlice) <= 1 {
			if command != "ADD" && command != "SAVE" && command != "DELETE" && command != "VIEW ALL" && command != "VIEW COMPLETED" && command != "VIEW PENDING" && command != "MARK COMPLETE" {
				fmt.Println("Invalid input, Try again")
				continue
			}
		}
		if len(inputSlice) >= 2 {
			arg = inputSlice[1]
			arg = strings.TrimSpace(arg)
		}

		switch command {
		case "ADD":
			addTask(arg)
		case "DELETE":
			argInt, _ := strconv.Atoi(strings.TrimSpace(arg))
			deleteTask(argInt)
		case "VIEW ALL":
			viewTasks()
		case "VIEW COMPLETED":
			viewCompletedTasks()
		case "VIEW PENDING":
			viewIncompletedTasks()
		case "MARK COMPLETE":
			argInt, _ := strconv.Atoi(strings.TrimSpace(arg))
			taskComplete(argInt)
		case "SAVE":
			save()
		default:
			fmt.Println("Command not supported")
		}
	}
}

func loadTasks() {
	jsonData, _ := os.ReadFile("tasks.json")
	json.Unmarshal(jsonData, &tasks)
	fmt.Println("\nSaved tasks successfully loaded to memeory")
}

func save() {
	jsonData, _ := json.MarshalIndent(tasks, "", " ")
	os.WriteFile("tasks.json", jsonData, 0644)
	fmt.Println("Changes Saved Successfully!")
}

func addTask(description string) {
	id := len(tasks) + 1
	for _, task := range tasks {
		if task.ID >= id {
			id = task.ID + 1
		}
	}
	task := Task{
		ID:          id,
		Description: description,
		Status:      false,
	}
	tasks = append(tasks, task)
	fmt.Println("Task added successfully")
	save()
}

func viewTasks() {
	fmt.Println("\nViewing all tasks")
	fmt.Print("ID\t|Status\t|\tDescription\n")
	if len(tasks) > 0 {
		for _, task := range tasks {
			status := "❌"
			if task.Status {
				status = "✅"
			}
			fmt.Printf("%v\t|%v\t|%v\n", task.ID, status, task.Description)
		}
	} else {
		fmt.Println("\tNo tasks for now")
	}
}

func viewCompletedTasks() {
	fmt.Println("\nViewing completed tasks")
	fmt.Print("ID\t|Status\t|\tDescription\n")
	if len(tasks) > 0 {
		found := false
		for _, task := range tasks {
			if task.Status == true {
				found = true
				status := "❌"
				if task.Status {
					status = "✅"
				}
				fmt.Printf("%v\t|%v\t|%v\n", task.ID, status, task.Description)
			}
		}
		if !found {
			fmt.Println("\tNo completed tasks yet")
		}
	} else {
		fmt.Println("\tNo tasks for now")
	}
}

func viewIncompletedTasks() {
	fmt.Println("\nViewing incompleted tasks")
	fmt.Print("ID\t|Status\t|\tDescription\n")
	if len(tasks) > 0 {
		found := false
		for _, task := range tasks {
			if task.Status == false {
				found = true
				status := "❌"
				if task.Status {
					status = "✅"
				}
				fmt.Printf("%v\t|%v\t|%v\n", task.ID, status, task.Description)
			}
		}
		if !found {
			fmt.Println("\tNo completed tasks yet")
		}
	} else {
		fmt.Println("\tNo tasks for now")
	}
}

func taskComplete(id int) {
	found := false
	if len(tasks) == 0 {
		fmt.Println("\tNo tasks here yet")
	} else {
		for i, task := range tasks {
			if task.ID == id {
				tasks[i].Status = true
				fmt.Printf("%v Marked completed\n", task.Description)
				found = true
				break
			}
		}
		if !found {
			fmt.Println("\tTask not found")
		}
	}
	save()
}

func deleteTask(id int) {
	found := false
	if len(tasks) == 0 {
		fmt.Printf("\tTask with ID:%v not found\n", id)
	} else {
		for i, task := range tasks {
			if task.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				fmt.Printf("Task %v Successfully Deleted\n", id)
				found = true
				break
			}
		}
		if !found {
			fmt.Println("\tTask not found")
		}
	}
	save()

}
