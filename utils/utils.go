package utils

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   error
}

type ErrorResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

const (
	UPDATE_FRIENDS       = "UPDATE_FRIENDS"       // action cannot be performed
	UPDATE_CLIQUES       = "UPDATE_CLIQUES"       // action cannot be performed
	UPDATE_CONVERSATIONS = "UPDATE_CONVERSATIONS" // action cannot be performed
)

type MessageData struct {
	MsgType string `json:"type"`
	Topic   string `json:"topic"`
	UserId  string `json:"userId"`
}
