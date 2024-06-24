package configurator

type CmdWant struct {
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

type AutomatorConfig struct {
	Headless bool `mapstructure:"headless"`
}

type GlobalConfig struct {
	Automator AutomatorConfig `mapstructure:"automator"`
}
