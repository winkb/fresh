package runner

import (
	"testing"
)

func TestIsWatchedFile(t *testing.T) {
	tests := []struct {
		file     string
		expected bool
	}{
		{"test.go", true},
		{"test.tpl", true},
		{"test.tmpl", true},
		{"test.html", true},
		{"test.css", false},
		{"test-executable", false},
		{"./tmp/test.go", false},
	}

	var s = NewMySetting(map[string]string{}, nil)

	for _, test := range tests {
		actual := isWatchedFile(s, test.file)

		if actual != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}

func TestShouldRebuild(t *testing.T) {
	tests := []struct {
		eventName string
		expected  bool
	}{
		{`"test.go": MODIFIED`, true},
		{`"test.tpl": MODIFIED`, false},
		{`"test.tmpl": DELETED`, false},
		{`"unknown.extension": DELETED`, true},
		{`"no_extension": ADDED`, true},
		{`"./a/path/test.go": MODIFIED`, true},
	}

	var s = NewMySetting(map[string]string{}, nil)

	for _, test := range tests {
		actual := shouldRebuild(s, test.eventName)

		if actual != test.expected {
			t.Errorf("Expected %v, got %v (event was '%s')", test.expected, actual, test.eventName)
		}
	}
}

func TestIsIgnoredFolder(t *testing.T) {
	tests := []struct {
		dir      string
		expected bool
	}{
		{"assets/node_modules", true},
		{"tmp/pid", true},
		{"app/controllers", false},
	}

	var s = NewMySetting(map[string]string{}, nil)

	for _, test := range tests {
		actual := isIgnoredFolder(s, test.dir)
		if actual != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}
