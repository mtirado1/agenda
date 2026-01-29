package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)



func main() {
	if len(os.Args) < 2 {
		fmt.Println("No filename provided.")
		os.Exit(1)
	}

	filename := os.Args[1]
	data, err := os.ReadFile(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	text := string(data)

	tasks := ParseTasks(text)

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
		fmt.Println("Today:")
		for _, task := range todayTasks {
			fmt.Printf("%s%s\n", strings.Repeat(" ", 14), task.ShortString())
		}
	} else {
		fmt.Println("No tasks for today.")
	}
	fmt.Println()

	if len(tomorrowTasks) > 0 {
		fmt.Println("Tomorrow:")
		for _, task := range tomorrowTasks {
			fmt.Printf("%s%s\n", strings.Repeat(" ", 14), task.ShortString())
		}
	} else {
		fmt.Println("No tasks for tomorrow.")
	}
	fmt.Println()

	if len(upcomingTasks) > 0 {
		fmt.Println("Upcoming:")
		for _, task := range upcomingTasks {
			fmt.Printf("  %s\n", task)
		}
	} else {
		fmt.Println("No upcoming tasks.")
	}
}
