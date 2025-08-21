package entities_test

import (
	. "scheduler/internal/entities"
	. "scheduler/test"
	"testing"

	"github.com/google/uuid"
)

func TestError_SetOption(t *testing.T) {
	error := NewError(uuid.NewString(), "task", "Invalid email ID", uuid.NewString(), make(map[string]string))

	Equals(t, 0, len(error.GetOptions()))

	error.SetOption("refund", "false")

	Equals(t, 1, len(error.GetOptions()))
}

func TestError_RemoveOption(t *testing.T) {
	error := NewError(uuid.NewString(), "task", "Invalid email ID", uuid.NewString(), make(map[string]string))

	Equals(t, 0, len(error.GetOptions()))

	error.SetOption("refund", "false")

	Equals(t, 1, len(error.GetOptions()))

	error.RemoveOption("refund")

	Equals(t, 0, len(error.GetOptions()))
}
