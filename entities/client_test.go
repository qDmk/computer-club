package entities

import (
	"testing"
)

func TestNewClientName(t *testing.T) {
	valid := []string{"dima", "dima123", "123dima123", "_-yadr0-_"}
	for _, name := range valid {
		_, err := NewClientName(name)
		if err != nil {
			t.Errorf("No error expected for %q", name)
		}
	}

	edges := []string{"a", "z", "0", "9", "-", "_"}
	for _, name := range edges {
		_, err := NewClientName(name)
		if err != nil {
			t.Errorf("No error expected for %q", name)
		}
	}

	invalid := []string{"", " ", "A", "Dima", "dima$", "konfeto4k@", "DROP TABLE Users"}
	for _, name := range invalid {
		_, err := NewClientName(name)
		if err == nil {
			t.Errorf("Error expected for %q", name)
		}
	}
}

func FuzzNewClientName(f *testing.F) {
	isValidChar := func(r rune) bool {
		return 'a' <= r && r <= 'z' ||
			'0' <= r && r <= '9' ||
			r == '-' ||
			r == '_'
	}

	isValidName := func(s string) bool {
		if s == "" {
			return false
		}

		for _, c := range s {
			if !isValidChar(c) {
				return false
			}
		}
		return true
	}

	f.Fuzz(func(t *testing.T, name string) {
		expected := isValidName(name)

		_, err := NewClientName(name)
		actual := err == nil

		if expected != actual {
			t.Errorf("For name %q expected to be valid: %v, actual: %v", name, expected, actual)
		}
	})
}
