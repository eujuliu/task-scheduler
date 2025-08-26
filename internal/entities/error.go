package entities

const (
	TypeErrorTransaction string = "TRANSACTION"
	TypeErrorTask        string = "TASK"
)

type Error struct {
	BaseEntity
	referenceId string
	userId      string
	kind        string
	reason      string
	options     map[string]string
}

func NewError(
	referenceId string,
	kind string,
	reason string,
	userId string,
	options map[string]string,
) *Error {
	return &Error{
		BaseEntity:  *NewBaseEntity(),
		referenceId: referenceId,
		kind:        kind,
		reason:      reason,
		options:     options,
		userId:      userId,
	}
}

func (e *Error) GetReferenceId() string {
	return e.referenceId
}

func (e *Error) GetUserId() string {
	return e.userId
}

func (e *Error) GetType() string {
	return e.kind
}

func (e *Error) GetReason() string {
	return e.reason
}

func (e *Error) SetOption(key string, value string) {
	e.options[key] = value
}

func (e *Error) RemoveOption(key string) {
	delete(e.options, key)
}

func (e *Error) GetOptions() map[string]string {
	return e.options
}
