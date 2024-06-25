package automator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"smuggr.xyz/spammr/common/configurator"

	"github.com/jaswdr/faker/v2"
)

var Faker = faker.New()

func generateRandomPhoneNumber() string {
	return fmt.Sprintf("%03d%03d%03d", rand.Intn(1000), rand.Intn(1000), rand.Intn(1000))
}

func getReplacementValue(placeholderName string, cmdFlags *configurator.CmdFlags, placeholderValues map[string]string) string {
	if value, exists := placeholderValues[placeholderName]; exists {
		return value
	}

	person := Faker.Person()
	internet := Faker.Internet()
	address := Faker.Address()

	var generatedValue string
	switch placeholderName {
	case "name":
		if cmdFlags.Want.Name != "" {
			generatedValue = cmdFlags.Want.Name
		} else {
			Logger.Warn("name not specified, generating random value")
			generatedValue = fmt.Sprintf("%s %s", person.FirstName(), person.LastName())
		}
	case "first_name":
		if cmdFlags.Want.FirstName != "" {
			generatedValue = cmdFlags.Want.FirstName
		} else {
			Logger.Warn("first name not specified, generating random value")
			generatedValue = person.FirstName()
		}
	case "last_name":
		if cmdFlags.Want.LastName != "" {
			generatedValue = cmdFlags.Want.LastName
		} else {
			Logger.Warn("last name not specified, generating random value")
			generatedValue = person.LastName()
		}
	case "phone_number":
		if cmdFlags.Want.PhoneNumber != "" {
			generatedValue = cmdFlags.Want.PhoneNumber
		} else {
			Logger.Warn("phone number not specified, generating random value")
			generatedValue = generateRandomPhoneNumber()
		}
	case "address":
		if cmdFlags.Want.Address != "" {
			generatedValue = cmdFlags.Want.Address
		} else {
			Logger.Warn("address not specified, generating random value")

			generatedValue = address.Address()
		}
	case "email":
		if cmdFlags.Want.Email != "" {
			generatedValue = cmdFlags.Want.Email
		} else {
			Logger.Warn("email not specified, generating random value")
			generatedValue = internet.Email()
		}
	default:
		generatedValue = ""
	}

	placeholderValues[placeholderName] = generatedValue

	return generatedValue
}

func replaceInString(value string, placeholderPattern *regexp.Regexp, cmdFlags *configurator.CmdFlags, placeholderValues map[string]string) string {
	matches := placeholderPattern.FindAllStringSubmatch(value, -1)
	for _, match := range matches {
		placeholder := match[0]
		placeholderName := match[1]
		replacement := getReplacementValue(placeholderName, cmdFlags, placeholderValues)
		value = strings.ReplaceAll(value, placeholder, replacement)
	}
	return value
}

func processActions(actions []Action, cmdFlags *configurator.CmdFlags, placeholderValues map[string]string) {
	placeholderPattern := regexp.MustCompile(`\{\{(.*?)\}\}`)

	for i := range actions {
		action := &actions[i]

		if actionValue, ok := action.Value.(string); ok {
			action.Value = replaceInString(actionValue, placeholderPattern, cmdFlags, placeholderValues)
		}

		for j := range action.OnFailure {
			onFailureAction := &action.OnFailure[j]
			processActions([]Action{*onFailureAction}, cmdFlags, placeholderValues)
		}
	}
}

func ReplacePlaceholders(automator *Automator, cmdFlags *configurator.CmdFlags) {
	placeholderValues := make(map[string]string)
	processActions(automator.Actions, cmdFlags, placeholderValues)

	automatorJSON, err := json.Marshal(automator)
	if err != nil {
		Logger.Error("Error marshalling automator:", err)
		return
	}

	automatorJSONString := string(automatorJSON)
	automatorJSONString = replaceInString(automatorJSONString, regexp.MustCompile(`\{\{(.*?)\}\}`), cmdFlags, placeholderValues)

	err = json.Unmarshal([]byte(automatorJSONString), automator)
	if err != nil {
		Logger.Error("Error unmarshalling modified automator JSON:", err)
	}
}