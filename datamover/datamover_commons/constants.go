package datamover_commons

import "errors"

const (
	AppName    = "Data Mover"
	AppVersion = "0.1.2"

	ModeProduction = "production"
	ModeDebug      = "debug"
)

// ---------------------------------------------------------------------------------------------------------------------
//		e v e n t s
// ---------------------------------------------------------------------------------------------------------------------

const (
	// EventOnDoStop application is stopping
	EventOnDoStop = "on_do_stop"

	// EventOnNextJobRun run another job in chain
	EventOnNextJobRun = "on_next_job_run"
)

// ---------------------------------------------------------------------------------------------------------------------
//		w o r k s p a c e s
// ---------------------------------------------------------------------------------------------------------------------

const (
	WpDirStart = "start"
	WpDirApp   = "app"
	WpDirWork  = "*"
)

// ---------------------------------------------------------------------------------------------------------------------
//		e r r o r s
// ---------------------------------------------------------------------------------------------------------------------

var (
	PanicSystemError                = errors.New("panic_system_error")
	JobAlreadyRunningError          = errors.New("job_already_running_error")
	DatabaseNotSupportedError       = errors.New("database_not_supported_error")
	ActionInvalidConfigurationError = errors.New("action_invalid_configuration_error")
)
