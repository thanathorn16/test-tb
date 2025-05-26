package register

type RequestBody struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type httpResponseTemplate struct {
	Result Result      `json:"result"`
	Data   interface{} `json:"data,omitempty"`
}

type Result struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
