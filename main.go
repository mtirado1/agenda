package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var boldStyle = lipgloss.NewStyle().Bold(true)

func main() {
	tasks := []Task{}

	for i := 1; i < len(os.Args); i += 1 {
		filename := os.Args[i]
		data, err := os.ReadFile(filename)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		text := string(data)

		fileTasks, err := ParseTasks(text)
		if err != nil {
			fmt.Println(fmt.Errorf("%s - %w", filename, err))
		} else {
			tasks = append(tasks, fileTasks...)
		}
	}

	slices.SortFunc(tasks, func(a, b Task) int {
		return a.date.Compare(b.date)
	})

	todayTasks := []Task{}
	tomorrowTasks := []Task{}
	upcomingTasks := []Task{}

	now := time.Now()
	for _, task := range tasks {
		if daysDifference(now, task) == 0 {
			todayTasks = append(todayTasks, task)
		} else if daysDifference(now, task) == 1 {
			tomorrowTasks = append(tomorrowTasks, task)
		} else if daysDifference(now, task) > 1 && len(upcomingTasks) < 5 {
			upcomingTasks = append(upcomingTasks, task)
		}
	}

	if len(todayTasks) > 0 {
		fmt.Println(boldStyle.Render("Today:"))
		for _, task := range todayTasks {
			fmt.Printf("%s%s\n", strings.Repeat(" ", 14), task.ShortString())
		}
	} else {
		fmt.Println("No tasks for today.")
	}
	fmt.Println()

	if len(tomorrowTasks) > 0 {
		fmt.Println(boldStyle.Render("Tomorrow:"))
		for _, task := range tomorrowTasks {
			fmt.Printf("%s%s\n", strings.Repeat(" ", 14), task.ShortString())
		}
	} else {
		fmt.Println("No tasks for tomorrow.")
	}
	fmt.Println()

	if len(upcomingTasks) > 0 {
		fmt.Println(boldStyle.Render("Upcoming:"))
		for _, task := range upcomingTasks {
			fmt.Printf("  %s\n", task)
		}
	} else {
		fmt.Println("No upcoming tasks.")
	}
}
