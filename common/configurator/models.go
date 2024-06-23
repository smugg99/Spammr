package configurator

type CmdWant struct {
	Boundary    string
	SessionID   string
	Name        string
	FirstName   string
	LastName    string
	PhoneNumber string
	Address     string
	Email       string
}

type CmdFlags struct {
	IsVerbose   bool
	RequestsDir string
	Want 		CmdWant
}

type GlobalConfig struct {

}
