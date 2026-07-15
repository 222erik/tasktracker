package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	Todo = iota
	InProgress
	Done
)

type Task struct {
	ID        int       `json:"id"`
	Desc      string    `json:"description"`
	Stauts    int       `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func TasksToJsonFile(tasks []Task, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	jsonData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func TasksFromJsonFile(tasks *[]Task, filename string) error {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileData, tasks)
	return err // Nil if everything goes good
}

func FindTaskByID(ID int, tasks []Task) (int, error) { // Returns the task index in the tasks slice
	for i, t := range tasks {
		if t.ID == ID {
			return i, nil
		}
	}
	return 0, fmt.Errorf("can't find task with ID %d", ID)
}

func main() {
	tasks := []Task{}
	err := TasksFromJsonFile(&tasks, "tasks.json")
	if err != nil {
		fmt.Println(err)
		// No need to exit
	}
	defer func() {
		err := TasksToJsonFile(tasks, "tasks.json")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	ID := func() int {
		IDUsed := func(ID int) bool {
			for _, t := range tasks {
				if ID == t.ID {
					return true
				}
			}
			return false
		}

		IDCount := 1
		for IDUsed(IDCount) { // While IDCount is used, increment IDCount
			IDCount++
		}
		return IDCount
	}

	if len(os.Args) == 1 {
		fmt.Println("You need to specify what to do")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) == 2 {
			fmt.Println("You need to give a description on what the task is")
			os.Exit(1)
		}

		desc := strings.Join(os.Args[2:], " ")

		tasks = append(tasks, Task{ID: ID(), Desc: desc, Stauts: Todo, CreatedAt: time.Now(), UpdatedAt: time.Now()})

	case "update":
		if len(os.Args) <= 3 {
			fmt.Println("You need to give a description of the updated task")
			os.Exit(1)
		}

		desc := strings.Join(os.Args[3:], " ")

		taskID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		taskIdx, err := FindTaskByID(taskID, tasks)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		taskToUpdate := &tasks[taskIdx]
		taskToUpdate.Desc = desc
		taskToUpdate.UpdatedAt = time.Now()

	case "delete":
		if len(os.Args) != 3 {
			fmt.Println("Write the ID of the task")
			os.Exit(1)
		}

		taskID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		taskIdx, err := FindTaskByID(taskID, tasks)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		tasks = append(tasks[:taskIdx], tasks[taskIdx+1:]...)

	case "mark-in-progress":
		if len(os.Args) != 3 {
			fmt.Println("Write the ID of the task")
			os.Exit(1)
		}

		taskID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		taskIdx, err := FindTaskByID(taskID, tasks)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		tasks[taskIdx].Stauts = InProgress

	case "mark-done":
		if len(os.Args) != 3 {
			fmt.Println("Write the ID of the task")
			os.Exit(1)
		}

		taskID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		taskIdx, err := FindTaskByID(taskID, tasks)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		tasks[taskIdx].Stauts = Done

	case "list":
		if len(os.Args) == 2 {
			for _, t := range tasks {
				fmt.Println(`"` + t.Desc + `"`)
				fmt.Println("ID:", t.ID)

				fmt.Printf("Stauts: ")
				switch t.Stauts {
				case Todo:
					fmt.Println("Todo")
				case InProgress:
					fmt.Println("In Progress")
				case Done:
					fmt.Println("Done")
				default:
					fmt.Println("Something went wrong")
					os.Exit(1)
				}

				fmt.Println("Created:", t.CreatedAt.Format("15:04"))
				fmt.Println("Last Updated:", t.UpdatedAt.Format("15:04"))
				fmt.Println() // Empty line
			}
		}

		if len(os.Args) == 3 {
			switch os.Args[2] {
			case "todo":
				for _, t := range tasks {
					if t.Stauts != Todo {
						continue
					}
					fmt.Println(`"` + t.Desc + `"`)
					fmt.Println("ID:", t.ID)
					fmt.Println("Created:", t.CreatedAt.Format("15:04"))
					fmt.Println("Last Updated:", t.UpdatedAt.Format("15:04"))
					fmt.Println() // Empty line
				}

			case "done":
				for _, t := range tasks {
					if t.Stauts != Done {
						continue
					}
					fmt.Println(`"` + t.Desc + `"`)
					fmt.Println("ID:", t.ID)
					fmt.Println("Created:", t.CreatedAt.Format("15:04"))
					fmt.Println("Last Updated:", t.UpdatedAt.Format("15:04"))
					fmt.Println() // Empty line
				}

			case "in-progress":
				for _, t := range tasks {
					if t.Stauts != InProgress {
						continue
					}
					fmt.Println(`"` + t.Desc + `"`)
					fmt.Println("ID:", t.ID)
					fmt.Println("Created:", t.CreatedAt.Format("15:04"))
					fmt.Println("Last Updated:", t.UpdatedAt.Format("15:04"))
					fmt.Println() // Empty line
				}
			default:
				fmt.Println("Not valid")
				os.Exit(1)
			}
		}
	default:
		fmt.Println("Not valid")
		os.Exit(1)
	}
}
