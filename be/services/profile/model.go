package profile

type httpResponseTemplate struct {
	Result Result      `json:"result"`
	Data   interface{} `json:"data,omitempty"`
}

type Result struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
