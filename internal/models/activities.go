package models

type Activity interface {
	Execute(task Task) error
}
