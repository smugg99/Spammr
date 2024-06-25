package logger

type Resource string

const (
	ResourceAutomator Resource = "AUTOMATOR"
	ResourceConfig    Resource = "CONFIG"
	ResourceEnv       Resource = "ENV"
)

// Resource messages
var (
	MsgResourceRegisterSuccess     = NewMessageWrapper("MsgResourceRegisterSuccess", "resource %s (%s) successfully registered", InfoLevel)
	MsgResourceFetchSuccess        = NewMessageWrapper("MsgResourceFetchSuccess", "resource %s (%s) successfully fetched", InfoLevel)
	MsgResourceUpdateSuccess       = NewMessageWrapper("MsgResourceUpdateSuccess", "resource %s (%s) successfully updated", InfoLevel)
	MsgResourceRemoveSuccess       = NewMessageWrapper("MsgResourceRemoveSuccess", "resource %s (%s) successfully removed", InfoLevel)
	MsgResourceAuthenticateSuccess = NewMessageWrapper("MsgResourceAuthenticateSuccess", "resource %s (%s) successfully authenticated", InfoLevel)
)

// Uncategorizated messages
var (
	MsgLoadingResource = NewMessageWrapper("MsgLoadingResource", "loading resource %s (%s)", InfoLevel)
	MsgResourceLoaded  = NewMessageWrapper("MsgResourceLoaded", "resource %s (%s) loaded", InfoLevel)
	MsgInitializing    = NewMessageWrapper("MsgInitializing", "initializing", InfoLevel)
	MsgCleaningUp      = NewMessageWrapper("MsgCleaningUp", "cleaning up", InfoLevel)
	MsgInitialized     = NewMessageWrapper("MsgInitialized", "initialized", InfoLevel)
	MsgCleanedUp       = NewMessageWrapper("MsgCleanedUp", "cleaned up", InfoLevel)
)

// Database errors
var (
	ErrResourceAlreadyExists = NewMessageWrapper("ErrResourceAlreadyExists", "resource %s (%s) already exists", ErrorLevel)
	ErrResourceNotFound      = NewMessageWrapper("ErrResourceNotFound", "resource %s (%s) not found", ErrorLevel)
	ErrInitializingResource    = NewMessageWrapper("ErrInitializingResource", "error initializing resource %s (%s)", ErrorLevel)
)

// Resource errors
var (
	ErrReadingResource         = NewMessageWrapper("ErrReadingResource", "error reading %s (%s)", FatalLevel)
	ErrFormattingResource      = NewMessageWrapper("ErrFormattingResource", "error formatting %s (%s)", FatalLevel)
	ErrHashingResource         = NewMessageWrapper("ErrHashingResource", "error hashing %s (%s)", FatalLevel)
	ErrResourcesDirectoryEmpty = NewMessageWrapper("ErrResourcesDirectoryEmpty", "resources (%s) directory is empty", FatalLevel)
)


// General errors
var (
	ErrUnexpected 		      = NewMessageWrapper("ErrUnexpected", "unexpected error: %s", ErrorLevel)
)

// Uncategorized errors
var (
	ErrInitializing               = NewMessageWrapper("ErrInitializing", "error initializing: %s", ErrorLevel)
	ErrCleaningUp                 = NewMessageWrapper("ErrCleaningUp", "error cleaning up: %s", ErrorLevel)
	ErrEnvVariableNotSet          = NewMessageWrapper("ErrEnvVariableNotSet", "environment variable not set: %s", ErrorLevel)
	ErrUnknownActionType          = NewMessageWrapper("ErrUnknownActionType", "unknown action type: %s", ErrorLevel)
	ErrActionReturnedFalse        = NewMessageWrapper("ErrActionReturnedFalse", "action returned false", ErrorLevel)
	ErrActionFailed		          = NewMessageWrapper("ErrActionFailed", "action %s failed: %v", ErrorLevel)
	ErrUnsupportedActionValueType = NewMessageWrapper("ErrUnsupportedActionValueType", "unsupported return action value type", ErrorLevel)
)
