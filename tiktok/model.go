package tiktok


type ErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	LogId   string `json:"log_id"`
}