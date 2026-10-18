package middleware

import "testing"

func TestRedactPassword(t *testing.T) {
	got := redact([]byte(`{"email":"a@b.c","password":"x"}`))
	if got == "" {
		t.Fatal("empty")
	}
	if contains(got, `"x"`) {
		t.Fatal(got)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
