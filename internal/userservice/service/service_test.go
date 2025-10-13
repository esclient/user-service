package service

import (
    "testing"
)

func TestGenerateVerificationCode(t *testing.T) {
    code, err := generateVerificationCode()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(code) != 6 {
        t.Fatalf("expected 6-digit code, got %q (len=%d)", code, len(code))
    }
    for _, r := range code {
        if r < '0' || r > '9' {
            t.Fatalf("expected numeric code, got %q", code)
        }
    }
}


