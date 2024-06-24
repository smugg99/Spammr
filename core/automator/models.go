package automator

type ActionType string
type Want string

const (
	ActionNavigate ActionType = "navigate"
	ActionWait     ActionType = "wait"
	ActionFill     ActionType = "fill"
)

const (
	WantName        Want = "name"
	WantFirstName   Want = "first_name"
	WantLastName    Want = "last_name"
	WantPhoneNumber Want = "phone_number"
	WantAddress     Want = "address"
	WantEmail       Want = "email"
)

type Action struct {
	Type     ActionType `json:"action"`
	Selector string     `json:"selector"`
	Duration int        `json:"duration"`
	Value    string     `json:"value"`
}

type Automator struct {
	Actions []Action `json:"actions"`
	Wants   []string `json:"wants"`
}