package entities

import "time"

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
}

func NewError(
	referenceId string,
	kind string,
	reason string,
	userId string,
) *Error {
	return &Error{
		BaseEntity:  *NewBaseEntity(),
		referenceId: referenceId,
		kind:        kind,
		reason:      reason,
		userId:      userId,
	}
}

func HydrateError(
	id, kind, userId, reason, referenceId string,
	createdAt, updatedAt time.Time,
) *Error {
	return &Error{
		BaseEntity: BaseEntity{
			id:        id,
			createdAt: createdAt,
			updatedAt: updatedAt,
		},

		kind:        kind,
		userId:      userId,
		reason:      reason,
		referenceId: referenceId,
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
