package main

import (
	"testing"
)

// TestHasHeadlessFlag tests the hasHeadlessFlag function with various inputs
func TestHasHeadlessFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "no arguments",
			args:     []string{},
			expected: false,
		},
		{
			name:     "only program name",
			args:     []string{"program"},
			expected: false,
		},
		{
			name:     "headless flag present",
			args:     []string{"program", "--headless"},
			expected: true,
		},
		{
			name:     "headless flag at beginning",
			args:     []string{"--headless", "program"},
			expected: true,
		},
		{
			name:     "headless flag in middle",
			args:     []string{"program", "--headless", "--other"},
			expected: true,
		},
		{
			name:     "headless flag at end",
			args:     []string{"program", "--other", "--headless"},
			expected: true,
		},
		{
			name:     "multiple headless flags",
			args:     []string{"program", "--headless", "--headless"},
			expected: true,
		},
		{
			name:     "similar but not exact flag - single dash",
			args:     []string{"program", "-headless"},
			expected: false,
		},
		{
			name:     "similar but not exact flag - with suffix",
			args:     []string{"program", "--headless-mode"},
			expected: false,
		},
		{
			name:     "similar but not exact flag - with prefix",
			args:     []string{"program", "--run-headless"},
			expected: false,
		},
		{
			name:     "similar but not exact flag - uppercase",
			args:     []string{"program", "--HEADLESS"},
			expected: false,
		},
		{
			name:     "similar but not exact flag - mixed case",
			args:     []string{"program", "--Headless"},
			expected: false,
		},
		{
			name:     "flag with equals sign",
			args:     []string{"program", "--headless=true"},
			expected: false,
		},
		{
			name:     "flag with space and value",
			args:     []string{"program", "--headless", "true"},
			expected: true,
		},
		{
			name:     "headless as part of another value",
			args:     []string{"program", "--mode", "headless"},
			expected: false,
		},
		{
			name:     "multiple different flags without headless",
			args:     []string{"program", "--verbose", "--debug", "--port", "8080"},
			expected: false,
		},
		{
			name:     "headless with other flags",
			args:     []string{"program", "--verbose", "--headless", "--debug"},
			expected: true,
		},
		{
			name:     "empty strings in args",
			args:     []string{"program", "", "--headless", ""},
			expected: true,
		},
		{
			name:     "only empty strings",
			args:     []string{"", "", ""},
			expected: false,
		},
		{
			name:     "whitespace variations",
			args:     []string{"program", " --headless", "--headless "},
			expected: false,
		},
		{
			name:     "headless with tabs",
			args:     []string{"program", "\t--headless"},
			expected: false,
		},
		{
			name:     "many arguments with headless",
			args:     []string{
				"program", "--config", "/path/to/config",
				"--verbose", "--debug", "--headless",
				"--port", "8080", "--host", "localhost",
			},
			expected: true,
		},
		{
			name:     "many arguments without headless",
			args:     []string{
				"program", "--config", "/path/to/config",
				"--verbose", "--debug", "--server",
				"--port", "8080", "--host", "localhost",
			},
			expected: false,
		},
		{
			name:     "nil slice",
			args:     nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasHeadlessFlag(tt.args)
			if result != tt.expected {
				t.Errorf("hasHeadlessFlag(%v) = %v, expected %v",
					tt.args, result, tt.expected)
			}
		})
	}
}

// TestHasHeadlessFlagEdgeCases tests additional edge cases
func TestHasHeadlessFlagEdgeCases(t *testing.T) {
	t.Run("very long argument list", func(t *testing.T) {
		args := make([]string, 1000)
		args[0] = "program"
		for i := 1; i < 999; i++ {
			args[i] = "--flag" + string(rune(i))
		}
		args[999] = "--headless"
		
		if !hasHeadlessFlag(args) {
			t.Error("should find --headless in very long argument list")
		}
	})

	t.Run("very long argument list without headless", func(t *testing.T) {
		args := make([]string, 1000)
		args[0] = "program"
		for i := 1; i < 1000; i++ {
			args[i] = "--flag" + string(rune(i))
		}
		
		if hasHeadlessFlag(args) {
			t.Error("should not find --headless in very long argument list without it")
		}
	})

	t.Run("repeated calls with same args", func(t *testing.T) {
		args := []string{"program", "--headless"}
		
		for i := 0; i < 10; i++ {
			if !hasHeadlessFlag(args) {
				t.Error("should consistently return true for same args")
			}
		}
	})

	t.Run("concurrent calls", func(t *testing.T) {
		args := []string{"program", "--headless"}
		
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				result := hasHeadlessFlag(args)
				if !result {
					t.Error("concurrent call should return true")
				}
				done <- true
			}()
		}
		
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

// TestHasHeadlessFlagDoesNotModifyInput verifies the function doesn't modify input
func TestHasHeadlessFlagDoesNotModifyInput(t *testing.T) {
	original := []string{"program", "--headless", "--other"}
	args := make([]string, len(original))
	copy(args, original)
	
	hasHeadlessFlag(args)
	
	for i := range args {
		if args[i] != original[i] {
			t.Errorf("hasHeadlessFlag modified input at index %d: got %q, expected %q",
				i, args[i], original[i])
		}
	}
}

// TestHasHeadlessFlagPerformance tests performance characteristics
func TestHasHeadlessFlagPerformance(t *testing.T) {
	t.Run("early termination when found", func(t *testing.T) {
		args := []string{"program", "--headless"}
		for i := 0; i < 1000; i++ {
			args = append(args, "--otherflag")
		}
		
		if !hasHeadlessFlag(args) {
			t.Error("should find --headless early and return")
		}
	})

	t.Run("full scan when not found", func(t *testing.T) {
		args := []string{"program"}
		for i := 0; i < 1000; i++ {
			args = append(args, "--otherflag")
		}
		
		if hasHeadlessFlag(args) {
			t.Error("should scan all args when --headless not present")
		}
	})
}

// BenchmarkHasHeadlessFlag benchmarks the function performance
func BenchmarkHasHeadlessFlag(b *testing.B) {
	benchmarks := []struct {
		name string
		args []string
	}{
		{
			name: "small_with_flag",
			args: []string{"program", "--headless"},
		},
		{
			name: "small_without_flag",
			args: []string{"program", "--other"},
		},
		{
			name: "medium_with_flag_at_start",
			args: []string{"program", "--headless", "--a", "--b", "--c", "--d", "--e"},
		},
		{
			name: "medium_with_flag_at_end",
			args: []string{"program", "--a", "--b", "--c", "--d", "--e", "--headless"},
		},
		{
			name: "medium_without_flag",
			args: []string{"program", "--a", "--b", "--c", "--d", "--e", "--f"},
		},
		{
			name: "large_with_flag_at_start",
			args: func() []string {
				args := []string{"program", "--headless"}
				for i := 0; i < 100; i++ {
					args = append(args, "--flag"+string(rune(i)))
				}
				return args
			}(),
		},
		{
			name: "large_without_flag",
			args: func() []string {
				args := []string{"program"}
				for i := 0; i < 100; i++ {
					args = append(args, "--flag"+string(rune(i)))
				}
				return args
			}(),
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				hasHeadlessFlag(bm.args)
			}
		})
	}
}

// TestHasHeadlessFlagWithRealWorldScenarios tests realistic usage patterns
func TestHasHeadlessFlagWithRealWorldScenarios(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
		scenario string
	}{
		{
			name:     "typical GUI launch",
			args:     []string{"/usr/bin/lylink-jellyfin"},
			expected: false,
			scenario: "User launches application via desktop icon",
		},
		{
			name:     "systemd service launch",
			args:     []string{"/usr/bin/lylink-jellyfin", "--headless"},
			expected: true,
			scenario: "Application started as systemd service",
		},
		{
			name:     "docker container launch",
			args:     []string{"./lylink-jellyfin", "--headless"},
			expected: true,
			scenario: "Application running in Docker container",
		},
		{
			name: "complex CLI with config",
			args: []string{
				"./lylink-jellyfin",
				"--config", "/etc/lylink/config.yaml",
				"--headless",
				"--log-level", "info",
			},
			expected: true,
			scenario: "Server deployment with custom config",
		},
		{
			name: "debug mode without headless",
			args: []string{
				"./lylink-jellyfin",
				"--debug",
				"--verbose",
				"--log-level", "debug",
			},
			expected: false,
			scenario: "Developer debugging with GUI",
		},
		{
			name:     "kubernetes pod",
			args:     []string{"/app/lylink-jellyfin", "--headless"},
			expected: true,
			scenario: "Running in Kubernetes pod",
		},
		{
			name: "CI/CD testing",
			args: []string{
				"./lylink-jellyfin",
				"--headless",
				"--test-mode",
			},
			expected: true,
			scenario: "Automated testing in CI/CD pipeline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasHeadlessFlag(tt.args)
			if result != tt.expected {
				t.Errorf("Scenario: %s\nhasHeadlessFlag(%v) = %v, expected %v",
					tt.scenario, tt.args, result, tt.expected)
			}
		})
	}
}

// TestHasHeadlessFlagBoundaryConditions tests boundary conditions
func TestHasHeadlessFlagBoundaryConditions(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "single element - exact match",
			args:     []string{"--headless"},
			expected: true,
		},
		{
			name:     "single element - no match",
			args:     []string{"--other"},
			expected: false,
		},
		{
			name:     "two elements - first match",
			args:     []string{"--headless", "arg2"},
			expected: true,
		},
		{
			name:     "two elements - second match",
			args:     []string{"arg1", "--headless"},
			expected: true,
		},
		{
			name:     "two elements - no match",
			args:     []string{"arg1", "arg2"},
			expected: false,
		},
		{
			name:     "empty string as only arg",
			args:     []string{""},
			expected: false,
		},
		{
			name:     "headless with null bytes",
			args:     []string{"program", "--headless\x00"},
			expected: false,
		},
		{
			name:     "headless with newline",
			args:     []string{"program", "--headless\n"},
			expected: false,
		},
		{
			name:     "headless with carriage return",
			args:     []string{"program", "--headless\r"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasHeadlessFlag(tt.args)
			if result != tt.expected {
				t.Errorf("hasHeadlessFlag(%v) = %v, expected %v",
					tt.args, result, tt.expected)
			}
		})
	}
}