package main

import (
	"testing"
	"time"
)

func TestParseMultipleTasks(t *testing.T) {
	tasks, err := ParseTasks(`

* 2026-02-01 6:00-9:00 Task description
* 2026-02-02           Task without time
* 2011-04-23 19:00     Just the start time

* Not a task
# Comment`)

	expectedTasks := []Task{
		Task{
			date:         time.Date(2026, 2, 1, 6, 0, 0, 0, time.UTC),
			hasStartTime: true,
			duration:     3 * time.Hour,
			description:  "Task description",
		},
		Task{
			date:         time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC),
			hasStartTime: false,
			description:  "Task without time",
		},
		Task{
			date:         time.Date(2011, 4, 23, 19, 0, 0, 0, time.UTC),
			hasStartTime: true,
			description:  "Just the start time",
		},
	}

	if err != nil {
		t.Errorf("Found parsing error.")
	}

	if len(tasks) != len(expectedTasks) {
		t.Errorf("Parsed %q tasks, expected %q", len(tasks), len(expectedTasks))
	}

	for i, task := range tasks {
		if task != expectedTasks[i] {
			t.Errorf("Parsed %q expected %q", task, expectedTasks[i])
		}
	}
}
