package settings

import "errors"

var (
	// ErrInvalidThreshold is returned when the threshold value is invalid
	ErrInvalidThreshold = errors.New("low margin threshold must be >= 0")
	
	// ErrSettingsNotFound is returned when shop settings are not found
	ErrSettingsNotFound = errors.New("shop settings not found")
)
