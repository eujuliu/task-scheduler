package queue

import "encoding/json"

type Event struct {
	ClientID string `json:"clientId"`
	Type     string `json:"type"`
	Data     string `json:"data"`
}

func NewEvent(id, kind string, data map[string]any) (*Event, error) {
	encoded, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	event := &Event{
		ClientID: id,
		Type:     kind,
		Data:     string(encoded),
	}

	return event, nil
}
