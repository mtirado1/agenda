package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var ErrInvalidRange = errors.New("Invalid time range.")

var (
	datePattern = regexp.MustCompile(`^[*\-+]\s+(\d{4}-\d\d-\d\d)\s+`)
	timePattern = regexp.MustCompile(`^(\d\d?:\d\d)(\s*-\s*(\d\d?:\d\d)|)`)
)

type Task struct {
	date         time.Time
	hasStartTime bool
	duration     time.Duration
	description  string
}

func (t Task) String() string {
	return fmt.Sprintf("%s  %s", t.date.Format("2006-01-02"), t.ShortString())
}

func (t Task) ShortString() string {
	if t.duration > 0*time.Minute {
		return fmt.Sprintf("%s-%s  %s", t.date.Format("15:04"), t.date.Add(t.duration).Format("15:04"), t.description)
	} else if !t.hasStartTime {
		return strings.Repeat(" ", 13) + t.description
	}
	return fmt.Sprintf("%s        %s", t.date.Format("15:04"), t.description)
}

func getDurationFromTimeRange(timeA, timeB string) (time.Duration, error) {
	tA, errA := time.Parse("15:04", timeA)
	tB, errB := time.Parse("15:04", timeB)

	if errA == nil && errB == nil {
		duration := tB.Sub(tA)
		if duration < 0 {
			return duration, ErrInvalidRange
		}
		return duration, nil
	}

	return time.Duration(0), ErrInvalidRange
}

func daysDifference(now time.Time, task Task) int {
	normalizedNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	normalizedTaskDate := time.Date(task.date.Year(), task.date.Month(), task.date.Day(), 0, 0, 0, 0, time.UTC)

	delta := normalizedTaskDate.Sub(normalizedNow)
	return int(delta.Hours()) / 24
}

func ParseTasks(text string) ([]Task, error) {
	lines := strings.Split(text, "\n")
	tasks := []Task{}

	for i, line := range lines {
		task := Task{}

		if match := datePattern.FindStringSubmatch(line); match != nil {
			date, err := time.Parse("2006-01-02", match[1])
			if err != nil {
				continue
			}
			task.date = date
			line = line[len(match[0]):]
		} else {
			continue
		}

		if match := timePattern.FindStringSubmatch(line); match != nil {
			t, err := time.Parse("15:04", match[1])
			if err != nil {
				return nil, fmt.Errorf("Line %d: %w", i+1, err)
			}

			task.hasStartTime = true
			task.date = time.Date(task.date.Year(), task.date.Month(), task.date.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
			duration, err := getDurationFromTimeRange(match[1], match[3])
			if err == nil {
				task.duration = duration
			} else if !errors.Is(err, ErrInvalidRange) {
				return nil, fmt.Errorf("Line %d: %w", i+1, err)
			}
			line = line[len(match[0]):]
		}

		task.description = strings.TrimSpace(line)
		if len(task.description) > 0 {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}
