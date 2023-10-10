package server

import (
	"testing"

	"github.com/EricNeid/go-webserver-gin/internal/verify"
)

func TestWrite(t *testing.T) {
	// arrange
	unit := LogService{Max: 99}
	// action
	unit.Write([]byte("Test-1"))
	unit.Write([]byte("Test-2"))
	unit.Write([]byte("Test-3"))
	// verify
	verify.Equals(t, 3, len(unit.Logs))
	verify.Equals(t, "Test-1", unit.Logs[0])
	verify.Equals(t, "Test-2", unit.Logs[1])
	verify.Equals(t, "Test-3", unit.Logs[2])
}

func TestWrite_withMax(t *testing.T) {
	// arrange
	unit := LogService{Max: 2}
	// action
	unit.Write([]byte("Test-1"))
	unit.Write([]byte("Test-2"))
	unit.Write([]byte("Test-3"))
	// verify
	verify.Equals(t, 2, len(unit.Logs))
	verify.Equals(t, "Test-2", unit.Logs[0])
	verify.Equals(t, "Test-3", unit.Logs[1])
}

func TestWrite_withLongMessage_shouldBeTruncated(t *testing.T) {
	// arrange
	unit := LogService{MaxMessageLength: 5}
	// action
	unit.Write([]byte("1234567890"))
	// verify
	verify.Equals(t, 1, len(unit.Logs))
	verify.Equals(t, "12345", unit.Logs[0])
}

func TestWrite_withShortMessage_shouldNotBeTruncated(t *testing.T) {
	// arrange
	unit := LogService{MaxMessageLength: 10}
	// action
	unit.Write([]byte("12345"))
	// verify
	verify.Equals(t, 1, len(unit.Logs))
	verify.Equals(t, "12345", unit.Logs[0])
}
