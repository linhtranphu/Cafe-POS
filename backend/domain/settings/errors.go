package settings

import "errors"

var (
	// ErrInvalidThreshold is returned when the threshold value is invalid
	ErrInvalidThreshold = errors.New("low margin threshold must be >= 0")
	
	// ErrSettingsNotFound is returned when shop settings are not found
	ErrSettingsNotFound = errors.New("shop settings not found")
	
	// ErrInvalidPaperWidth is returned when paper width is not 58 or 80
	ErrInvalidPaperWidth = errors.New("paper width must be 58 or 80 mm")
	
	// ErrInvalidLabelSize is returned when label size is not valid
	ErrInvalidLabelSize = errors.New("label size must be one of: 40x30, 50x30, 60x40")
)
