package common

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
)

func GenerateProcessorId() string {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown-host"
	}
	return fmt.Sprintf("%s-%d-%v", host, os.Getpid(), strings.ReplaceAll(uuid.New().String(), "-", ""))
}
