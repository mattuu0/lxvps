package structs

type HttpResult struct {
	Code    int      `json:"code"`
	Message string `json:"message"`
	Error   error    `json:"error"`
	Success bool     `json:"success"`
}
