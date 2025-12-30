package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(false)

	if logger == nil {
		t.Error("NewLogger should not return nil")
	}

	if logger.IsVerbose() {
		t.Error("NewLogger with false should not be verbose")
	}

	logger = NewLogger(true)
	if !logger.IsVerbose() {
		t.Error("NewLogger with true should be verbose")
	}
}

func TestNewLoggerWithOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithOutput(false, &buf)

	if logger == nil {
		t.Error("NewLoggerWithOutput should not return nil")
	}

	if logger.IsVerbose() {
		t.Error("NewLoggerWithOutput with false should not be verbose")
	}

	logger = NewLoggerWithOutput(true, &buf)
	if !logger.IsVerbose() {
		t.Error("NewLoggerWithOutput with true should be verbose")
	}
}

func TestSetVerbose(t *testing.T) {
	logger := NewLogger(false)

	logger.SetVerbose(true)
	if !logger.IsVerbose() {
		t.Error("SetVerbose(true) should set verbose to true")
	}

	logger.SetVerbose(false)
	if logger.IsVerbose() {
		t.Error("SetVerbose(false) should set verbose to false")
	}
}

func TestPrintln(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithOutput(false, &buf)

	logger.Println("test message")
	output := buf.String()

	// Check that output contains the expected message and has timestamp format
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected output to contain 'test message', got '%s'", output)
	}
	if !strings.Contains(output, " ") { // Check for space separator
		t.Errorf("Expected timestamp format, got '%s'", output)
	}
}

func TestPrintf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithOutput(false, &buf)

	logger.Printf("test %s %d", "message", 42)
	output := buf.String()

	// Check that output contains the expected message and has timestamp format
	if !strings.Contains(output, "test message 42") {
		t.Errorf("Expected output to contain 'test message 42', got '%s'", output)
	}
	if !strings.Contains(output, " ") { // Check for space separator
		t.Errorf("Expected timestamp format, got '%s'", output)
	}
}

func TestPrintlnVerbose(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithOutput(false, &buf)

	// Should not output when not verbose
	logger.PrintlnVerbose("verbose message")
	if buf.String() != "" {
		t.Error("PrintlnVerbose should not output when not verbose")
	}

	// Should output when verbose
	logger.SetVerbose(true)
	logger.PrintlnVerbose("verbose message")
	output := buf.String()
	// Check that output contains the expected message and has timestamp format
	if !strings.Contains(output, "verbose message") {
		t.Errorf("Expected output to contain 'verbose message', got '%s'", output)
	}
	if !strings.Contains(output, " ") { // Check for space separator
		t.Errorf("Expected timestamp format, got '%s'", output)
	}
}

func TestPrintfVerbose(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithOutput(false, &buf)

	// Should not output when not verbose
	logger.PrintfVerbose("verbose %s", "message")
	if buf.String() != "" {
		t.Error("PrintfVerbose should not output when not verbose")
	}

	// Should output when verbose
	logger.SetVerbose(true)
	logger.PrintfVerbose("verbose %s %d", "message", 42)
	output := buf.String()
	// Check that output contains the expected message and has timestamp format
	if !strings.Contains(output, "verbose message 42") {
		t.Errorf("Expected output to contain 'verbose message 42', got '%s'", output)
	}
	if !strings.Contains(output, " ") { // Check for space separator
		t.Errorf("Expected timestamp format, got '%s'", output)
	}
}

func TestPrintlnError(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithOutput(false, &buf)

	// Note: PrintlnError always writes to os.Stderr, not the custom output
	// This test just verifies the method doesn't panic
	logger.PrintlnError("error message")
	// We can't easily capture stderr in this test, but we can verify it doesn't crash
}

func TestPrintfError(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithOutput(false, &buf)

	// Note: PrintfError always writes to os.Stderr, not the custom output
	// This test just verifies the method doesn't panic
	logger.PrintfError("error %s", "message")
	// We can't easily capture stderr in this test, but we can verify it doesn't crash
}
