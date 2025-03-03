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
	Priority    int    `json:"priority"`
	DueDate     string `json:"dueDate"`
}

type Taskinator interface {
	AddTask(description string, priority int, dueDate string)
	ViewTasks()
	ViewCompletedTasks()
	ViewIncompletedTasks()
	TaskComplete(id int)
	Save()
	DeleteTask(id int)
}

type TaskList struct {
	tasks []Task
}

var myList TaskList

func main() {
	/*
	 ADD-Go to church on sunday-3-2025/03/03
	 ADD-Read a book-2-2025/03/09
	 ADD-Read the bible-1-2025/03/07
	 ADD-Drink water-0-2025/03/02
	 ADD-Go to the Gym at 15:00-3-2025/03/04
	*/

	myList.LoadTasks()

	commands := map[string]string{
		"Add new task":           "ADD-\"Task Description\"-Priority(0 to 3)-Due Date(YYYY//MM/DD)",
		"Delete a task":          "DELETE-\"Task ID\"",
		"View all tasks":         "VIEW ALL",
		"View Completed tasks":   "VIEW COMPLETED",
		"View Pending tasks":     "VIEW PENDING",
		"Mark task as Completed": "MARK COMPLETE-\"Task ID\"",
	}

	for {
		fmt.Println("\n\tCLI TASKINATOR IS RUNNING")
		fmt.Print("\nTo perform a command, Please enter any of these prompts below (e.g. ADD-Read a book-2-2025/03/03) \n\n")
		for command, prompt := range commands {
			fmt.Printf("To %v:: \t%v\n", command, prompt)
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nEnter a command: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		inputSlice := strings.Split(input, "-")
		command := strings.ToUpper(strings.TrimSpace(inputSlice[0]))
		var (
			arg1 string
			arg2 int
			arg3 string
		)
		if len(inputSlice) <= 1 {
			if command != "ADD" && command != "SAVE" && command != "DELETE" && command != "VIEW ALL" && command != "VIEW COMPLETED" && command != "VIEW PENDING" && command != "MARK COMPLETE" {
				fmt.Println("Invalid input, Try again")
				continue
			}
		} else if len(inputSlice) >= 2 {
			arg1 = strings.TrimSpace(inputSlice[1])
			arg2 = 3
			arg3 = "2024/02/28"
		}
		if len(inputSlice) >= 3 {
			arg2, _ = strconv.Atoi(strings.TrimSpace(inputSlice[2]))
		}
		if len(inputSlice) >= 4 {
			arg3 = strings.TrimSpace(inputSlice[3])
		}

		switch command {
		case "ADD":
			fmt.Println(inputSlice)
			fmt.Printf("command:%v\ndesc:%v\nprio:%v\ndate:%v\n", command, arg1, arg2, arg3)
			myList.AddTask(arg1, arg2, arg3)
		case "VIEW ALL":
			myList.ViewTasks()
		case "VIEW COMPLETED":
			myList.ViewCompletedTasks()
		case "VIEW PENDING":
			myList.ViewIncompletedTasks()
		case "MARK COMPLETE":
			arg1Int, _ := strconv.Atoi(arg1)
			myList.TaskComplete(arg1Int)
		case "SAVE":
			myList.Save()
		case "DELETE":
			arg1Int, _ := strconv.Atoi(arg1)
			myList.DeleteTask(arg1Int)
		default:
			fmt.Println("Command not supported")
		}
	}
}

func (tList *TaskList) LoadTasks() {
	jsonData, _ := os.ReadFile("tasks.json")
	json.Unmarshal(jsonData, &tList.tasks)
	fmt.Println("\nSaved tasks successfully loaded to memeory")
}

func (tList *TaskList) Save() {
	jsonData, _ := json.MarshalIndent(tList.tasks, "", " ")
	os.WriteFile("tasks.json", jsonData, 0644)
	fmt.Println("Changes Saved Successfully!")
}

func (tList *TaskList) AddTask(description string, priority int, dueDate string) {
	id := len(tList.tasks) + 1
	for _, task := range tList.tasks {
		if task.ID >= id {
			id = task.ID + 1
		}
	}
	task := Task{
		ID:          id,
		Description: description,
		Status:      false,
		Priority:    priority,
		DueDate:     dueDate,
	}
	tList.tasks = append(tList.tasks, task)
	fmt.Println("Task added successfully")
	myList.Save()
}

func (tList *TaskList) ViewTasks() {
	fmt.Println("\nViewing all tasks")
	fmt.Print("ID\t|Status\t|\tDescription\t| Priority | Due Date\n")
	if len(tList.tasks) > 0 {
		for _, task := range tList.tasks {
			status := "❌"
			if task.Status {
				status = "✅"
			}
			fmt.Printf("%v\t|%v\t|%v\t|%v\t|%v\n", task.ID, status, task.Description, task.Priority, task.DueDate)
		}
	} else {
		fmt.Println("\tNo tasks for now")
	}
}

func (tList *TaskList) ViewCompletedTasks() {
	fmt.Println("\nViewing completed tasks")
	fmt.Print("ID\t|Status\t|\tDescription\n")
	if len(tList.tasks) > 0 {
		found := false
		for _, task := range tList.tasks {
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

func (tList *TaskList) ViewIncompletedTasks() {
	fmt.Println("\nViewing incompleted tasks")
	fmt.Print("ID\t|Status\t|\tDescription\n")
	if len(tList.tasks) > 0 {
		found := false
		for _, task := range tList.tasks {
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

func (tList *TaskList) TaskComplete(id int) {
	found := false
	if len(tList.tasks) == 0 {
		fmt.Println("\tNo tasks here yet")
	} else {
		for i, task := range tList.tasks {
			if task.ID == id {
				tList.tasks[i].Status = true
				fmt.Printf("%v Marked completed\n", task.Description)
				found = true
				break
			}
		}
		if !found {
			fmt.Println("\tTask not found")
		}
	}
	myList.Save()
}

func (tList *TaskList) DeleteTask(id int) {
	found := false
	if len(tList.tasks) == 0 {
		fmt.Printf("\tTask with ID:%v not found\n", id)
	} else {
		for i, task := range tList.tasks {
			if task.ID == id {
				myList.tasks = append(tList.tasks[:i], tList.tasks[i+1:]...)
				fmt.Printf("Task %v Successfully Deleted\n", id)
				found = true
				break
			}
		}
		if !found {
			fmt.Println("\tTask not found")
		}
	}
	myList.Save()

}
