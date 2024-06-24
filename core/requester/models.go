package requester

type Want string

const (
	WantHash        Want = "hash"
	WantBoundary    Want = "boundary"
	WantSessionID   Want = "session_id"
	WantName        Want = "name"
	WantFirstName   Want = "first_name"
	WantLastName    Want = "last_name"
	WantPhoneNumber Want = "phone_number"
	WantAddress     Want = "address"
	WantEmail       Want = "email"
)

type RequestTemplate struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Wants   []Want            `json:"wants"`
}

type ResponseDetails struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

type ResponseLog struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}