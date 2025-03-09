package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	Priority    int       `json:"priority"`
	DueDate     time.Time `json:"dueDate"`
}

type Taskinator interface {
	AddTask(description string, priority int, dueDate time.Time)
	Save()
	LoadTasks(ch chan bool)
	ViewTasks() string
	ViewCompletedTasks() string
	ViewIncompletedTasks() string
	TaskComplete(id int) string
	DeleteTask(id int) string
}

type TaskList struct {
	tasks []Task
}

type TaskRequest struct {
	Command     string
	Description string
	ID          int
	Priority    int
	DueDate     time.Time
}

var (
	myList Taskinator

	taskQueue = make(chan TaskRequest, 5)
	results   = make(chan string)
	wg        sync.WaitGroup
	mu        sync.Mutex
)

func main() {
	myList = &TaskList{}
	/*
		echo -e "ADD-Test1-1-2025/03/06\nADD-Test2-2-2025/03/07\nMARK-1\nDELETE-2\nALL" | xargs -I {} -P 5 go run main.go "{}"
		ADD-Go to church on sunday-3-2025/03/03
		ADD-Read a book-2-2025/03/09
		ADD-Read the bible-1-2025/03/07
		ADD-Drink water-0-2025/03/02
		ADD-Go to the Gym at 15:00-3-2025/03/04
		ADD-Read Golang book-2-2025/03/05
		MARK-1
		DELETE-2
		ADD-Complete project-3-2025/03/10
		ALL
	*/
	ch := make(chan bool)
	ch2 := make(chan string)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go Quit(signalChan)

	log.Println("Loading Tasks...")
	go myList.LoadTasks(ch)
	if <-ch {
		log.Println("Done!")
	}

	commands := map[string]string{
		"Add new task":           "ADD-\"Task Description\"-Priority(0 to 3)-Due Date(YYYY//MM/DD)",
		"Delete a task":          "DELETE-\"Task ID\"",
		"View all tasks":         "ALL",
		"View Completed tasks":   "DONE",
		"View Pending tasks":     "PENDING",
		"Mark task as Completed": "MARK-\"Task ID\"",
	}

	for w := 1; w <= 3; w++ {
		go Worker(w, taskQueue, results)
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

		// go func(input string) {
		command, arg1, id, priority, dDate, err := ParseInput(input)
		if err != nil {
			fmt.Println("Error:", err)
			// return
			continue
		}
		// mu.Lock()
		fmt.Println(<-ch2)
		// mu.Unlock()
		wg.Add(1)
		taskQueue <- TaskRequest{Command: command, Description: arg1, ID: id, Priority: priority, DueDate: dDate}
		fmt.Println(<-results)

		// go ProcessTask(command, arg1, id, priority, dDate, ch2)
		// fmt.Println(<-ch2)
		// }(input)
	}
}

func ParseInput(userInput string) (command string, arg1 string, id int, priority int, dDate time.Time, err error) {
	userInput = strings.TrimSpace(userInput)
	inputSlice := strings.Split(userInput, "-")
	command = strings.ToUpper(strings.TrimSpace(inputSlice[0]))

	if len(inputSlice) < 2 {
		if slices.Contains([]string{"SAVE", "ALL", "DONE", "PENDING"}, command) {
			return
		}
		return "", "", 0, 0, time.Time{}, fmt.Errorf("invalid command format")
	}

	switch command {
	case "ADD":
		if len(inputSlice) != 4 {
			return "", "", 0, 0, time.Time{}, fmt.Errorf("ADD requires description, priority, and due date")
		}
		arg1 = strings.TrimSpace(inputSlice[1])
		priority, err = strconv.Atoi(strings.TrimSpace(inputSlice[2]))
		if err != nil {
			return "", "", 0, 0, time.Time{}, fmt.Errorf("invalid priority value")
		}
		dDate, err = time.Parse("2006/01/02", strings.TrimSpace(inputSlice[3]))
		if err != nil {
			return "", "", 0, 0, time.Time{}, fmt.Errorf("invalid date format, use YYYY/MM/DD")
		}
	case "DELETE", "MARK":
		id, err = strconv.Atoi(strings.TrimSpace(inputSlice[1]))
		if err != nil {
			return "", "", 0, 0, time.Time{}, fmt.Errorf("invalid ID")
		}
	default:
		return "", "", 0, 0, time.Time{}, fmt.Errorf("unsupported command")
	}
	return
}

func ProcessTask(command string, arg1 string, id int, priority int, dDate time.Time, ch2 chan string) {
	// Add task to the wait group
	// Send task ID to taskQueue
	// Create a unique response channel per request
	// Send/return the processed result to ch2

	wg.Add(1)
	responseChan := make(chan string)
	taskQueue <- TaskRequest{Command: command, Description: arg1, ID: id, Priority: priority, DueDate: dDate}
	// fmt.Printf("Command: %v\nDescription: %v\nID: %v\nPriority: %v\nDue Date: %v\n", command, arg1, id, priority, dDate)
	go func() {
		result := <-responseChan
		ch2 <- result
		results <- result
		close(responseChan)
	}()
}

func (tList *TaskList) LoadTasks(ch chan bool) {
	defer close(ch)

	jsonData, err := os.ReadFile("tasks.json")
	if err != nil {
		log.Printf("Error reading tasks file: %v", err)
		ch <- false
		return
	}
	err = json.Unmarshal(jsonData, &tList.tasks)
	if err != nil {
		log.Printf("Error unmarshalling tasks: %v", err)
		ch <- false
		return
	}
	ch <- true
}

func (tList *TaskList) Save() {
	jsonData, _ := json.MarshalIndent(tList.tasks, "", " ")
	os.WriteFile("tasks.json", jsonData, 0644)
	fmt.Println("File Updated!")
}

func (tList *TaskList) AddTask(description string, priority int, dueDate time.Time) {
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
	// fmt.Printf("Command: \nDescription: %v\nID: \nPriority: %v\nDue Date: %v\n", description, priority, dueDate)

	if description != "" {
		tList.tasks = append(tList.tasks, task)
		fmt.Println("Task added successfully")
		myList.Save()
	} else {
		fmt.Println("Empty description")
	}
}

func (tList *TaskList) ViewTasks() string {
	var result strings.Builder
	result.WriteString("\nViewing all tasks\n")
	result.WriteString("ID\t|Status\t|\tDescription\t| Priority | Due Date\n")
	if len(tList.tasks) > 0 {
		for _, task := range tList.tasks {
			status := "❌"
			if task.Status {
				status = "✅"
			}
			result.WriteString(fmt.Sprintf("%v\t|%v\t|%v\t|%v\t|%v\n", task.ID, status, task.Description, task.Priority, task.DueDate))
		}
	} else {
		result.WriteString("\tNo tasks for now\n")
	}
	return result.String()
}

func (tList *TaskList) ViewCompletedTasks() string {
	var result strings.Builder
	result.WriteString("\nViewing completed tasks")
	result.WriteString("ID\t|Status\t|\tDescription\n")
	if len(tList.tasks) > 0 {
		found := false
		for _, task := range tList.tasks {
			if task.Status == true {
				found = true
				status := "❌"
				if task.Status {
					status = "✅"
				}
				result.WriteString(fmt.Sprintf("%v\t|%v\t|%v\n", task.ID, status, task.Description))
			}
		}
		if !found {
			result.WriteString("\tNo completed tasks yet")
		}
	} else {
		result.WriteString("\tNo tasks for now")
	}
	return result.String()

}

func (tList *TaskList) ViewIncompletedTasks() string {
	var result strings.Builder
	result.WriteString("\nViewing incompleted tasks")
	result.WriteString("ID\t|Status\t|\tDescription\n")
	if len(tList.tasks) > 0 {
		found := false
		for _, task := range tList.tasks {
			if task.Status == false {
				found = true
				status := "❌"
				if task.Status {
					status = "✅"
				}
				result.WriteString(fmt.Sprintf("%v\t|%v\t|%v\n", task.ID, status, task.Description))
			}
		}
		if !found {
			result.WriteString("\tNo completed tasks yet")
		}
	} else {
		result.WriteString("\tNo tasks for now")
	}
	return result.String()
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

func Worker(id int, tasks <-chan TaskRequest, results chan<- string) {
	var result string
	for task := range tasks {
		// fmt.Printf("\tLOG FROM WORKER\nCommand: %v\nDescription: %v\nID: %v\nPriority: %v\nDue Date: %v\n", task.Command, task.Description, task.ID, task.Priority, task.DueDate)
		switch task.Command {
		case "ADD":
			myList.AddTask(task.Description, task.Priority, task.DueDate)
		case "ALL":
			result = myList.ViewTasks()
		case "DONE":
			result = myList.ViewCompletedTasks()
		case "PENDING":
			result = myList.ViewIncompletedTasks()
		case "SAVE":
			myList.Save()
			result = "Tasks saved successfully"
		case "MARK":
			result = myList.TaskComplete(task.ID)
		case "DELETE":
			result = myList.DeleteTask(task.ID)
		default:
			result = "Command not supported\nTry Again!"
		}
		results <- fmt.Sprintf("Worker %d: %s", id, result)
		wg.Done()
	}
}

func Quit(quit chan os.Signal) {
	<-quit
	fmt.Println("\nTidying things Up...")
	myList.Save()
	close(taskQueue)
	wg.Wait()
	fmt.Println("\nBis Bald!")
	os.Exit(0)
}
