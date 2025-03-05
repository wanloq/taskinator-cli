package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
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
	Save()
	LoadTasks(ch chan bool)
	ViewTasks()
	ViewCompletedTasks()
	ViewIncompletedTasks()
	TaskComplete(id int) string
	DeleteTask(id int) string
}

type TaskList struct {
	tasks []Task
}

var (
	myList  Taskinator
	command string
	arg1    string
	arg2    int
	arg3    string

	taskQueue = make(chan int, 5)
	results   = make(chan string)
	wg        sync.WaitGroup
)

func main() {
	myList = &TaskList{}
	/*
	 ADD-Go to church on sunday-3-2025/03/03
	 ADD-Read a book-2-2025/03/09
	 ADD-Read the bible-1-2025/03/07
	 ADD-Drink water-0-2025/03/02
	 ADD-Go to the Gym at 15:00-3-2025/03/04
	*/
	ch := make(chan bool)
	ch2 := make(chan string)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit)
	go Quit(quit)

	status := "Loading Tasks..."
	go myList.LoadTasks(ch)
	log.Println(status)
	if <-ch == true {
		status = "Done!"
	}
	log.Println(status)
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

		ParseInput(input)

		switch command {
		case "ADD":
			myList.AddTask(arg1, arg2, arg3)
		case "VIEW ALL":
			myList.ViewTasks()
		case "VIEW COMPLETED":
			myList.ViewCompletedTasks()
		case "VIEW PENDING":
			myList.ViewIncompletedTasks()
		case "SAVE":
			myList.Save()
		default:
			id, _ := strconv.Atoi(arg1)
			go ProcessTask(id, ch2)
			fmt.Println(<-ch2)
		}
	}
}

func Quit(quit chan os.Signal) {
	<-quit
	fmt.Println("\nTidying things Up...")
	myList.Save()
	// close(taskQueue)
	// wg.Wait()
	fmt.Println("\nBis Später!")
	os.Exit(0)
}

func ParseInput(userInput string) {
	userInput = strings.TrimSpace(userInput)
	inputSlice := strings.Split(userInput, "-")
	command = strings.ToUpper(strings.TrimSpace(inputSlice[0]))
	if len(inputSlice) <= 1 {
		if command != "ADD" && command != "SAVE" && command != "DELETE" && command != "VIEW ALL" && command != "VIEW COMPLETED" && command != "VIEW PENDING" && command != "MARK COMPLETE" {
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

}

func ProcessTask(id int, ch chan string) {
	switch command {
	case "MARK COMPLETE":
		msg := myList.TaskComplete(id)
		ch <- msg
	case "DELETE":
		msg := myList.DeleteTask(id)
		ch <- msg
	default:
		ch <- "Command not supported\nTry Again!"
	}
}

func (tList *TaskList) LoadTasks(ch chan bool) {
	defer close(ch)

	jsonData, _ := os.ReadFile("tasks.json")
	json.Unmarshal(jsonData, &tList.tasks)
	// fmt.Println("\nSaved tasks successfully loaded to memeory")
	ch <- true
}

func (tList *TaskList) Save() {
	jsonData, _ := json.MarshalIndent(tList.tasks, "", " ")
	os.WriteFile("tasks.json", jsonData, 0644)
	fmt.Println("File Updated!")
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

func (tList *TaskList) TaskComplete(id int) string {
	if len(tList.tasks) == 0 {
		return fmt.Sprintln("\tNo tasks here yet")
	}
	for i, task := range tList.tasks {
		if task.ID == id {
			tList.tasks[i].Status = true
			myList.Save()
			return fmt.Sprintf("%v Successfully Completed\n", task.Description)
		}
	}
	return fmt.Sprintf("\tTask with ID: %v not found\n", id)
}

func (tList *TaskList) DeleteTask(id int) string {
	if len(tList.tasks) == 0 {
		return fmt.Sprintln("\tNo tasks here yet")
	}
	for i, task := range tList.tasks {
		if task.ID == id {
			tList.tasks = append(tList.tasks[:i], tList.tasks[i+1:]...)
			myList.Save()
			return fmt.Sprintf("%v Successfully Removed\n", task.Description)
		}
	}
	return fmt.Sprintf("\tTask with ID: %v not found\n", id)
}

func Worker(id int, tasks <-chan int, results chan<- string) {
	for taskID := range tasks {
		result := myList.TaskComplete(taskID) // Mark task as complete
		results <- fmt.Sprintf("Worker %d: %s", id, result)
		wg.Done() // Signal that a task is done
	}
}
