package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Status int

const (
	Todo Status = iota
	InProgress
	Done
)

func (s Status) String() string {
	switch s {
	case Todo:
		return "Todo"
	case InProgress:
		return "In Progress"
	case Done:
		return "Done"
	default:
		return "Not Valid"
	}
}

type Task struct {
	ID        int       `json:"id"`
	Desc      string    `json:"description"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func TasksToJSONFile(tasks []Task, filename string) error {
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
	return err
}

func TasksFromJSONFile(tasks *[]Task, filename string) error {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileData, tasks)
	return err // Nil if everything goes good
}

func GetTaskIdx(tasks []Task) int {
	taskID, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	taskIdx, err := func() (int, error) {
		for i, t := range tasks {
			if t.ID == taskID {
				return i, nil
			}
		}
		return 0, fmt.Errorf("can't find task with ID %d", taskID)
	}()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return taskIdx
}

func (t Task) Print() {
	fmt.Println(`"` + t.Desc + `"`)
	fmt.Println("ID:", t.ID)
	fmt.Println("Created:", t.CreatedAt.Format("2006-01-02 15:04"))
	fmt.Println("Last Updated:", t.UpdatedAt.Format("2006-01-02 15:04"))
	fmt.Println() // Empty line
}

func main() {
	tasks := []Task{}
	err := TasksFromJSONFile(&tasks, "tasks.json")
	if err != nil && !os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		err := TasksToJSONFile(tasks, "tasks.json")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	ID := func() int {
		maxID := 0
		for _, t := range tasks {
			if t.ID > maxID {
				maxID = t.ID
			}
		}
		return maxID + 1
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

		tasks = append(tasks, Task{ID: ID(), Desc: desc, Status: Todo, CreatedAt: time.Now(), UpdatedAt: time.Now()})

	case "update":
		if len(os.Args) <= 3 {
			fmt.Println("You need to give a description of the updated task")
			os.Exit(1)
		}

		desc := strings.Join(os.Args[3:], " ")

		taskToUpdate := &tasks[GetTaskIdx(tasks)]
		taskToUpdate.Desc = desc
		taskToUpdate.UpdatedAt = time.Now()

	case "delete":
		if len(os.Args) != 3 {
			fmt.Println("Write the ID of the task")
			os.Exit(1)
		}

		taskIdx := GetTaskIdx(tasks)
		tasks = append(tasks[:taskIdx], tasks[taskIdx+1:]...)

	case "mark-in-progress":
		if len(os.Args) != 3 {
			fmt.Println("Write the ID of the task")
			os.Exit(1)
		}

		tasks[GetTaskIdx(tasks)].Status = InProgress

	case "mark-done":
		if len(os.Args) != 3 {
			fmt.Println("Write the ID of the task")
			os.Exit(1)
		}

		tasks[GetTaskIdx(tasks)].Status = Done

	case "list":
		if len(os.Args) == 2 {
			for _, t := range tasks {
				fmt.Println(`"` + t.Desc + `"`)
				fmt.Println("ID:", t.ID)
				fmt.Println("Status:", t.Status.String())
				fmt.Println("Created:", t.CreatedAt.Format("2006-01-02 15:04"))
				fmt.Println("Last Updated:", t.UpdatedAt.Format("2006-01-02 15:04"))
				fmt.Println() // Empty line
			}
		}

		if len(os.Args) == 3 {
			switch os.Args[2] {
			case "todo":
				for _, t := range tasks {
					if t.Status != Todo {
						continue
					}
					t.Print()
				}

			case "done":
				for _, t := range tasks {
					if t.Status != Done {
						continue
					}
					t.Print()
				}

			case "in-progress":
				for _, t := range tasks {
					if t.Status != InProgress {
						continue
					}
					t.Print()
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
