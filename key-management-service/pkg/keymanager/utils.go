package keymanager

import "strings"

// ToSecurityMeasures returns the base field size and subgroup order based on the security level
func ToSecurityMeasures(level string) (uint32, uint32) {
	level = strings.ToLower(level)
	switch level {
	case "low":
		return 128, 256 // Lower security level, suitable for testing or less critical applications with 128-bit base field and 256-bit subgroup order
	case "medium":
		return 160, 512 // Balanced level of security and performance with 160-bit base field and 512-bit subgroup order
	case "high":
		return 256, 1024 // Higher security level with 256-bit base field and 1024-bit subgroup order
	default:
		return 160, 512 // Default to Balanced level of security and performance with 160-bit base field and 512-bit subgroup order
	}
}
