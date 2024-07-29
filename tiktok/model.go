package tiktok

type ErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	LogId   string `json:"log_id"`
}


type AccessTokenManagement struct {
	AccessToken string `json:"access_token"`
	TokenType 	string `json:"token_type"`
	ExpiresIn 	int64 `json:"expires_in"`
}