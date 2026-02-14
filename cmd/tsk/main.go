package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/zarldev/claude-mngr/internal/task"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	path, err := task.DefaultPath()
	if err != nil {
		fatal(err)
	}
	s := task.NewStore(path)

	switch os.Args[1] {
	case "add":
		runAdd(s, os.Args[2:])
	case "list":
		runList(s, os.Args[2:])
	case "done":
		runDone(s, os.Args[2:])
	case "rm":
		runRm(s, os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func runAdd(s *task.Store, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: tsk add <title>")
		os.Exit(1)
	}
	title := strings.Join(args, " ")
	t, err := s.Add(title)
	if err != nil {
		fatal(err)
	}
	fmt.Printf("added task %d: %s\n", t.ID, t.Title)
}

func runList(s *task.Store, args []string) {
	var filter *bool
	for _, a := range args {
		switch a {
		case "--done":
			v := true
			filter = &v
		case "--pending":
			v := false
			filter = &v
		default:
			fmt.Fprintf(os.Stderr, "unknown flag: %s\n", a)
			os.Exit(1)
		}
	}

	tasks, err := s.List(filter)
	if err != nil {
		fatal(err)
	}
	if len(tasks) == 0 {
		fmt.Println("no tasks")
		return
	}
	for _, t := range tasks {
		check := "[ ]"
		if t.Done {
			check = "[x]"
		}
		fmt.Printf("%3d  %s  %-40s  %s\n", t.ID, check, t.Title, age(t.CreatedAt))
	}
}

func runDone(s *task.Store, args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: tsk done <id>")
		os.Exit(1)
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid id: %s\n", args[0])
		os.Exit(1)
	}
	if err := s.Done(id); err != nil {
		fatal(err)
	}
	fmt.Printf("task %d marked done\n", id)
}

func runRm(s *task.Store, args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: tsk rm <id>")
		os.Exit(1)
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid id: %s\n", args[0])
		os.Exit(1)
	}
	if err := s.Remove(id); err != nil {
		fatal(err)
	}
	fmt.Printf("task %d removed\n", id)
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: tsk <command> [args]

commands:
  add <title>       create a new task
  list [--done|--pending]  show tasks
  done <id>         mark a task as done
  rm <id>           remove a task`)
}

func age(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "tsk: %v\n", err)
	os.Exit(1)
}
