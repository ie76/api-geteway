package errors

import (
	"fmt"
	"log"
	"time"
)

const (
	ErrInvalidInput              = "ErrInvalidInput"
	ErrInvalidEnvInput           = "ErrInvalidEnvInput"
	ErrServiceNotFound           = "ErrServiceNotFound"
	ErrDatabaseError             = "ErrDatabaseError"
	ErrRedisConnect              = "ErrRedisConnect"
	ErrHash                      = "ErrHash"
	ErrUnauthorized              = "ErrUnauthorized"
	ErrPlanNotFound              = "ErrPlanNotFound"
	ErrEnvNotFound               = "ErrEnvNotFound"
	ErrConfigInit                = "ErrConfigInit"
	ErrUserNotFound              = "ErrUserNotFound"
	ErrCreateUser                = "ErrCreateUser"
	ErrTokenGeneration           = "ErrTokenGeneration"
	ErrTokenDecrypt              = "ErrTokenDecrypt"
	ErrInsufficientCredits       = "ErrInsufficientCredits"
	ErrCreatePlan                = "ErrCreatePlan"
	ErrServiceFailedToLoadConfig = "ErrServiceFailedToLoadConfig"
	ErrParseUrl                  = "ErrParseUrl"
	ErrServiceHTTPRequest        = "errServiceHTTPGetRequest"
	ErrServiceReadBody           = "ErrServiceReadBody"
	ErrCacheError                = "ErrCacheError"
)

type Error struct {
	Code      string
	Message   string
	Time      time.Time
	Retryable bool
	logger    *log.Logger
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error %s: %s time at %s", e.Code, e.Message, e.Time.Format(time.RFC3339))
}

func New(code string) *Error {
	message, retryable := errorMessages[code]
	return &Error{
		Code:      code,
		Message:   message.string,
		Time:      time.Now(),
		Retryable: retryable,
	}
}

func NewWithLogger(code string, logger *log.Logger) *Error {
	err := New(code)
	err.logger = logger
	return err
}

var errorMessages = map[string]struct {
	string
	bool
}{
	ErrInvalidInput:              {"Invalid input", false},
	ErrServiceNotFound:           {"Service not found", false},
	ErrInvalidEnvInput:           {"Invalid env value", false},
	ErrDatabaseError:             {"Database error", true},
	ErrRedisConnect:              {"Redis Database error", true},
	ErrHash:                      {"Hash not possible", true},
	ErrUnauthorized:              {"Unauthorized access", false},
	ErrPlanNotFound:              {"Plan not found", false},
	ErrTokenGeneration:           {"Token not generated", false},
	ErrTokenDecrypt:              {"Token couldn't be decrypted", false},
	ErrEnvNotFound:               {"Env file not found or doesn't exists", false},
	ErrConfigInit:                {"config can't be initiated", false},
	ErrUserNotFound:              {"User not found", false},
	ErrCreateUser:                {"can't create User", false},
	ErrInsufficientCredits:       {"Insuffisant credits", false},
	ErrCreatePlan:                {"Create plan failed", false},
	ErrServiceFailedToLoadConfig: {"Service Failed to load config", false},
	ErrParseUrl:                  {"Unable to parse url", false},
	ErrServiceHTTPRequest:        {"Service failed to execute HTTP Request", false},
	ErrCacheError:                {"Failed to cache data", false},
	ErrServiceReadBody:           {"Service can't parse body", false},
}

func RegisterErrorMessage(code, message string, retryable bool) {
	errorMessages[code] = struct {
		string
		bool
	}{message, retryable}
}

func (e *Error) Log() {
	if e.logger == nil {
		e.logger = log.Default()
	}
	e.logger.Printf("%s", e.Error())
}

func Retry(attempts int, sleep time.Duration, f func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = f()
		if err == nil {
			return nil
		}

		if customErr, ok := err.(*Error); ok && !customErr.Retryable {
			break
		}

		time.Sleep(sleep)
	}
	return err
}
