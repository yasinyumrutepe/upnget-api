package globals

type Response struct {
	Error   interface{} `json:"error,omitempty"`
	Body    interface{} `json:"body,omitempty"`
	Message string      `json:"message"`
}
