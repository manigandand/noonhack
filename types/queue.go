package types

// QueueInput ...
type QueueInput struct {
	ServiceName string      `json:"service_name"`
	Data        interface{} `json:"data"`
}
