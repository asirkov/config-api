package api

import (
	"fmt"
)

type ServiceStatus struct {
	Service string
	Failed  bool
	Message string
}

type WarningValidation struct {
	Msg string
}

func (e *WarningValidation) Error() string {
	return e.Msg
}

type InfoNotFound struct {
	Entity string
}

func (i *InfoNotFound) Error() string {
	if i.Entity == "" {
		return "was not found"
	}
	return fmt.Sprintf("%s was not found", i.Entity)
}

type ErrorNotModified struct {
}

func (e *ErrorNotModified) Error() string {
	return "The item was not modified"
}

type ErrorDuplicateEntries struct {
	Duplicates interface{}
}

func (e *ErrorDuplicateEntries) Error() string {
	return fmt.Sprintf("Duplicate entries found: %v", e.Duplicates)
}
