package queue

const (
	SEND_EMAIL_QUEUE       string = "email-task"
	GET_TASKS_RESULT_QUEUE string = "task-result"
	SEND_EVENTS_QUEUE      string = "events"

	SEND_EMAIL_KEY      string = "task.email.send"
	GET_TASK_RESULT_KEY string = "task.result"
	SEND_EVENTS_KEY     string = "events.send"

	TASK_EXCHANGE   string = "tasks"
	EVENTS_EXCHANGE string = "events"
)

var WorkersQueues = map[string]string{
	"email": SEND_EMAIL_KEY,
}
