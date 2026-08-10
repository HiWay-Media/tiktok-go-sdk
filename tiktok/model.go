package tiktok

// ErrorObject is the error envelope of the /v2/... endpoints. TikTok fills it
// on every response, successful or not: Code is "ok" when the call succeeded.
type ErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	LogId   string `json:"log_id"`
}

// IsError reports whether the envelope describes a failure. An empty code is
// treated as success: some endpoints omit it when nothing went wrong.
func (e ErrorObject) IsError() bool {
	return e.Code != "" && e.Code != "ok"
}

// OAuthErrorObject is the error envelope of /v2/oauth/token/, which does not
// use ErrorObject: it reports failures with a flat shape, and still answers
// HTTP 200. Modelled separately on purpose — folding both shapes into one
// struct makes every decode fail on whichever shape it did not expect.
type OAuthErrorObject struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	LogId            string `json:"log_id"`
}

// IsError reports whether the envelope describes a failure.
func (e OAuthErrorObject) IsError() bool { return e.Error != "" }

type AccessTokenManagement struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}
