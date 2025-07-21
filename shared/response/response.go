package response

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}