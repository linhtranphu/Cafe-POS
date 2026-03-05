package printing

import "encoding/base64"

// isBase64Content checks if a string is likely base64-encoded content
func isBase64Content(s string) bool {
	// Base64 strings should be reasonably long and contain only valid base64 characters
	if len(s) < 100 {
		return false
	}
	
	// Try to decode - if successful, it's likely base64
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
