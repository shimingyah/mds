package utils

import (
	"strings"

	"github.com/google/uuid"
)

// NewUUID uuid format not "-"
func NewUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

// RawUUID raw uuid format
func RawUUID() string {
	return uuid.New().String()
}
