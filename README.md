# taskinator-cli
As my first mini-project on this journey, I’ll build this simple CLI-based task manager (TASKINATOR) that allows users to:

✅ Add tasks (with a description and status)  
✅ List all tasks (showing task number, description, and completion status)  
✅ Mark tasks as complete  
✅ Delete tasks.  
✅ List completed tasks.  
✅ List incompleted tasks.  
✅ Use JSON encoding/decoding.  
✅ Save tasks to file.  
✅ Load tasks from file.  

## Additional Improvements

1.  Add a priority.  
2.  Add a due date.  

# Summary

During this mini-project, I applied core Go concepts to build a functional CLI Task Manager. Here a breakdown:

## V1 

✔ Achieved all the functionalities needed but did not implement any extra features to boost efficiency.  

## V2 
✔ Achieved all the needed functionalities.  
✔ Implemented an interface for handling the tasklist.  
✔ Implemented goroutines and channels for handling concurrent operations on tasks.  
✔ Implemnets a parser.  
✔ Implements a process task function for handling operations that can be done concurrently.  
✔ 
✔ 
✔ 

## 1. Go Syntax & Program Structure

✔ Structured my Go program with functions, loops, and conditionals to handle different commands efficiently.  
✔ Applied maps to store command prompts dynamically.  

## 2. Working with Data Structures

✔ Stored all tasks dynamically using slices.  
✔ Modeled task properties using struct (ID, Description, Status).  
✔ Mapped commands to their corresponding syntax with maps.  

## 3. File Handling, Persistence, Serialization & Deserialization (os & encoding/json)

✔ Reading & Writing Files with the os package.  
✔ Handling JSON Data with encoding/json package.  
✔ Creating a Persistent Task List.  

## 4. User Input Handling (bufio.Reader)

✔ Used bufio.NewReader(os.Stdin) to read user input.  
✔ Used strings.Split to parse commands.  
✔ Applied string manipulation (TrimSpace(), ToUpper()) to clean inputs.  

## 5. Error Handling & Validation

✔ Checked for invalid input cases.  
✔ Ensured valid task IDs when marking tasks as complete or deleting them.  
✔ Detailed error handling has not been properly implemented (INTENTIONAL).  

## 6. Functions & Code Organization

✔ Encapsulated logic in functions (e.g., addTask(), viewTasks()).  
✔ Used parameters to make functions reusable (func deleteTask(id int)).  

## 7. Iterating Over Data (for loops)

✔ Used 'for' to range slices and maps to achieve diffent desired effects.  
✔ Applied conditionals inside loops to filter tasks.  

## 8. Extra Challenges Implemented

✔ Dynamic Task ID generation (preventing duplicate IDs after deletion).  
✔ More User-Friendly Display (e.g. ✅ for 'Completed' and ❌ for Pending tasks instead of true/false).  
