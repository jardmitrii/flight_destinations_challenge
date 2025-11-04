package main

import (
	"fmt"
	"strings"
)

func validateAirportCode(origin string) error {
	if len(origin) != 3 {
		return fmt.Errorf("invalid origin %s: length != 3 ", origin)
	}
	if !isLetters(origin) {
		return fmt.Errorf("invalid origin %s: contains non-latin letters", origin)
	}

	return nil
}

func isLetters(s string) bool {
	for _, r := range strings.ToUpper(s) {
		if !(r >= 'A' && r <= 'Z') {
			return false
		}
	}
	return true
}

func getDestination(origin string, destination string) (string, bool) {
	destination = strings.ToUpper(destination)
	origin = strings.ToUpper(origin)
	if strings.HasPrefix(destination, origin) {
		return strings.TrimPrefix(destination, origin), true
	}
	if strings.HasSuffix(destination, origin) {
		return strings.TrimSuffix(destination, origin), true
	}

	return "", false
}
