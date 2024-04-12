package models

import (
	"time"
)

type Task struct {
	Id         int
	Name       string
	Desc       string
	Enddate    time.Time
	Priority   Priority
	Status     Status
	Categories []Category
}

type Priority struct {
	Id   int
	Name string
}

type Status struct {
	Id   int
	Name string
}

type Category struct {
	Id   int
	Name string
}
