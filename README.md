# Tasktracker
Simple task tracker CLI in Go

## Compile
```bash
go build main.go
```

## How to use 📜
- ```./main add <task>``` To add a task
- ```./main uptade <task-id> <new-task>``` To update task
- ```./main delete <task-id>``` To delete task
- ```./main mark-in-progress <task-id>``` To mark a task in progress
- ```./main mark-done <task-id>``` To mark a task done
- ```./main list``` To list all tasks
- ```./main list todo``` To list all tasks with status todo
- ```./main list in-progress``` To list all tasks in progress
- ```./main list done``` To list all tasks that are done

## Features
- Saves tasks to json file ✔️
- Keeps track on when tasks are created and uptaded 📆
- Each task can have status as Todo, In Progress, or Done 📝

## Task Tracker on roadmap.sh link 🖇️
https://roadmap.sh/projects/task-tracker
