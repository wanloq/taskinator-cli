package main

// import (
// 	"bufio"
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// type Task struct {
// 	ID          int    `json:"id"`
// 	Description string `json:"description"`
// 	Status      bool   `json:"status"`
// 	Priority    int    `json:"priority"`
// 	DueDate     string `json:"dueDate"`
// }

// var tasks []Task

// func main() {
// 	/*
// 	 ADD-Go to church on sunday-3-2025/03/03
// 	 ADD-Read a book-2-2025/03/09
// 	 ADD-Read the bible-1-2025/03/07
// 	 ADD-Drink water-0-2025/03/02
// 	 ADD-Go to the Gym at 15:00-3-2025/03/04
// 	*/
// 	loadTasks()
// 	commands := map[string]string{
// 		"Add new task":           "ADD-\"Task Description\"-Priority(0 to 3)-Due Date",
// 		"Delete a task":          "DELETE-\"Task ID\"",
// 		"View all tasks":         "VIEW ALL",
// 		"View Completed tasks":   "VIEW COMPLETED",
// 		"View Pending tasks":     "VIEW PENDING",
// 		"Mark task as Completed": "MARK COMPLETE-\"Task ID\"",
// 	}

// 	for {
// 		fmt.Println("\n\tCLI TASKINATOR IS RUNNING")
// 		fmt.Print("\nTo perform a command, Please enter any of these prompts below (e.g. ADD-Read a book-2-2025/03/03) \n\n")
// 		for command, prompt := range commands {
// 			fmt.Printf("To %v:: \t%v\n", command, prompt)
// 		}
// 		reader := bufio.NewReader(os.Stdin)
// 		fmt.Print("\nEnter a command: ")
// 		input, _ := reader.ReadString('\n')
// 		input = strings.TrimSpace(input)
// 		inputSlice := strings.Split(input, "-")
// 		command := strings.ToUpper(strings.TrimSpace(inputSlice[0]))
// 		var (
// 			arg1 string
// 			arg2 int
// 			arg3 string
// 		)
// 		if len(inputSlice) <= 1 {
// 			if command != "ADD" && command != "SAVE" && command != "DELETE" && command != "VIEW ALL" && command != "VIEW COMPLETED" && command != "VIEW PENDING" && command != "MARK COMPLETE" {
// 				fmt.Println("Invalid input, Try again")
// 				continue
// 			}
// 		} else if len(inputSlice) >= 2 {
// 			arg1 = strings.TrimSpace(inputSlice[1])
// 			arg2 = 3
// 			arg3 = "2024/02/28"
// 		}
// 		if len(inputSlice) >= 3 {
// 			arg2, _ = strconv.Atoi(strings.TrimSpace(inputSlice[2]))
// 		}
// 		if len(inputSlice) >= 4 {
// 			arg3 = strings.TrimSpace(inputSlice[3])
// 		}

// 		switch command {
// 		case "ADD":
// 			fmt.Println(inputSlice)
// 			fmt.Printf("command:%v\ndesc:%v\nprio:%v\ndate:%v\n", command, arg1, arg2, arg3)
// 			addTask(arg1, arg2, arg3)
// 		case "VIEW ALL":
// 			viewTasks()
// 		case "VIEW COMPLETED":
// 			viewCompletedTasks()
// 		case "VIEW PENDING":
// 			viewIncompletedTasks()
// 		case "MARK COMPLETE":
// 			arg1Int, _ := strconv.Atoi(arg1)
// 			taskComplete(arg1Int)
// 		case "SAVE":
// 			save()
// 		case "DELETE":
// 			arg1Int, _ := strconv.Atoi(arg1)
// 			deleteTask(arg1Int)
// 		default:
// 			fmt.Println("Command not supported")
// 		}
// 	}
// }

// func loadTasks() {
// 	jsonData, _ := os.ReadFile("tasks.json")
// 	json.Unmarshal(jsonData, &tasks)
// 	fmt.Println("\nSaved tasks successfully loaded to memeory")
// }

// func save() {
// 	jsonData, _ := json.MarshalIndent(tasks, "", " ")
// 	os.WriteFile("tasks.json", jsonData, 0644)
// 	fmt.Println("Changes Saved Successfully!")
// }

// func addTask(description string, priority int, dueDate string) {
// 	id := len(tasks) + 1
// 	for _, task := range tasks {
// 		if task.ID >= id {
// 			id = task.ID + 1
// 		}
// 	}
// 	task := Task{
// 		ID:          id,
// 		Description: description,
// 		Status:      false,
// 		Priority:    priority,
// 		DueDate:     dueDate,
// 	}
// 	tasks = append(tasks, task)
// 	fmt.Println("Task added successfully")
// 	save()
// }

// func viewTasks() {
// 	fmt.Println("\nViewing all tasks")
// 	fmt.Print("ID\t|Status\t|\tDescription\t| Priority | Due Date\n")
// 	if len(tasks) > 0 {
// 		for _, task := range tasks {
// 			status := "❌"
// 			if task.Status {
// 				status = "✅"
// 			}
// 			fmt.Printf("%v\t|%v\t|%v\t|%v\t|%v\n", task.ID, status, task.Description, task.Priority, task.DueDate)
// 		}
// 	} else {
// 		fmt.Println("\tNo tasks for now")
// 	}
// }

// func viewCompletedTasks() {
// 	fmt.Println("\nViewing completed tasks")
// 	fmt.Print("ID\t|Status\t|\tDescription\n")
// 	if len(tasks) > 0 {
// 		found := false
// 		for _, task := range tasks {
// 			if task.Status == true {
// 				found = true
// 				status := "❌"
// 				if task.Status {
// 					status = "✅"
// 				}
// 				fmt.Printf("%v\t|%v\t|%v\n", task.ID, status, task.Description)
// 			}
// 		}
// 		if !found {
// 			fmt.Println("\tNo completed tasks yet")
// 		}
// 	} else {
// 		fmt.Println("\tNo tasks for now")
// 	}
// }

// func viewIncompletedTasks() {
// 	fmt.Println("\nViewing incompleted tasks")
// 	fmt.Print("ID\t|Status\t|\tDescription\n")
// 	if len(tasks) > 0 {
// 		found := false
// 		for _, task := range tasks {
// 			if task.Status == false {
// 				found = true
// 				status := "❌"
// 				if task.Status {
// 					status = "✅"
// 				}
// 				fmt.Printf("%v\t|%v\t|%v\n", task.ID, status, task.Description)
// 			}
// 		}
// 		if !found {
// 			fmt.Println("\tNo completed tasks yet")
// 		}
// 	} else {
// 		fmt.Println("\tNo tasks for now")
// 	}
// }

// func taskComplete(id int) {
// 	found := false
// 	if len(tasks) == 0 {
// 		fmt.Println("\tNo tasks here yet")
// 	} else {
// 		for i, task := range tasks {
// 			if task.ID == id {
// 				tasks[i].Status = true
// 				fmt.Printf("%v Marked completed\n", task.Description)
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			fmt.Println("\tTask not found")
// 		}
// 	}
// 	save()
// }

// func deleteTask(id int) {
// 	found := false
// 	if len(tasks) == 0 {
// 		fmt.Printf("\tTask with ID:%v not found\n", id)
// 	} else {
// 		for i, task := range tasks {
// 			if task.ID == id {
// 				tasks = append(tasks[:i], tasks[i+1:]...)
// 				fmt.Printf("Task %v Successfully Deleted\n", id)
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			fmt.Println("\tTask not found")
// 		}
// 	}
// 	save()

// }
