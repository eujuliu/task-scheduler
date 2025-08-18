package services_test

// import (
// 	"scheduler/internal/entities"
// 	"scheduler/internal/errors"
// 	. "scheduler/test"
// 	"testing"
// 	"time"

// 	"github.com/google/uuid"
// )

// func TestCreateTaskService(t *testing.T) {
// 	teardown := Setup(t)
// 	defer teardown(t)

// 	user, _ := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

// 	user.SetCredits(100)

// 	task, err := CreateTaskService.Execute("email_send", user.GetId(), time.Now().AddDate(0, 1, 0), uuid.NewString())

// 	Ok(t, err)

// 	transaction, err = TransactionRepository.GetByReferenceId(task.GetId())
// 	user, _ = GetUserService.Execute("testuser", "Password@123")

// 	Ok(t, err)
// 	Equals(t, entities.TaskCost["email_send"], transaction.GetCost())
// 	Equals(t, "task_send", transaction.GetType())
// 	Equals(t, 100-entities.TaskCost["email_send"], user.GetCredits())
// }

// func TestCreateTaskService_UserNotFound(t *testing.T) {
// 	teardown := Setup(t)
// 	defer teardown(t)

// 	_, err := CreateTaskService.Execute("email_send", uuid.NewString(), time.Now().AddDate(0, 1, 0), uuid.NewString())

// 	Equals(t, errors.USER_NOT_FOUND_ERROR().Error(), err.Error())
// }

// func TestCreateTaskService_UserWithoutCredits(t *testing.T) {
// 	teardown := Setup(t)
// 	defer teardown(t)

// 	user, _ := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

// 	_, err := CreateTaskService.Execute("email_send", user.GetId(), time.Now().AddDate(0, 1, 0), uuid.NewString())

// 	Equals(t, errors.USER_WITHOUT_CREDITS().Error(), err.Error())
// }
