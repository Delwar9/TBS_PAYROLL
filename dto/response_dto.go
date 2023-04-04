package dto

type ResponseDto struct {
	IsSuccess  bool        `json:"isSuccess"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Count      int         `json:"count,omitempty"`
	Payload    interface{} `json:"payload"`
}
