package binary

import (
	"fmt"
	"strconv"
	"strings"
)

// parseTag parses the tag to extract length specification
func parseTag(tag string) (uint32, error) {
	if tag == "" {
		return 0, fmt.Errorf("empty tag")
	}

	// If tag is "-", it means to ignore the tag
	if tag == "-" {
		return 0, fmt.Errorf("ignore tag")
	}

	// Try to parse as integer
	if length, err := strconv.ParseUint(tag, 10, 32); err == nil {
		return uint32(length), nil
	}

	// Try to parse as "len:N" format
	if strings.HasPrefix(tag, "len:") {
		parts := strings.Split(tag, ":")
		if len(parts) == 2 {
			if length, err := strconv.ParseUint(parts[1], 10, 32); err == nil {
				return uint32(length), nil
			}
		}
	}

	return 0, fmt.Errorf("invalid tag format: %s", tag)
}