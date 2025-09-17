package queue

const (
	SEND_EMAIL_QUEUE string = "email-queue"
	SEND_EMAIL_KEY   string = "task.email.send"

	TASK_EXCHANGE string = "task-exchange"
)

var AvailableQueues = map[string]string{
	"email": SEND_EMAIL_KEY,
}
