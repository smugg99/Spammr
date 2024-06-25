package automator

type ActionType string
type Want string

const (
	ActionNavigate ActionType = "navigate"
	ActionWait     ActionType = "wait"
	ActionFill     ActionType = "fill"
	ActionReturn   ActionType = "return"
	ActionPrint    ActionType = "print"
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
	Type      ActionType  `json:"action"`
	Selector  string      `json:"selector,omitempty"`
	Duration  int         `json:"duration,omitempty"`
	Value     interface{} `json:"value,omitempty"`
	OnFailure []Action    `json:"onFailure,omitempty"`
}

type Automator struct {
	Name    string          `json:"name"`
	Actions []Action        `json:"actions"`
	Wants   []string        `json:"wants"`
	Has     map[string]Want `json:"-"`
}