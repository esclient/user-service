package repository

import (
    "testing"
    "time"
)

func TestIsCodeMismatch(t *testing.T) {
    if err := isCodeMismatch("654321", "123456"); err == nil {
        t.Fatalf("expected mismatch error, got nil")
    }
    if err := isCodeMismatch("123456", "123456"); err != nil {
        t.Fatalf("unexpected error for equal codes: %v", err)
    }
}

func TestIsCodeExpired(t *testing.T) {
    past := time.Now().Add(-VerificationCodeLifetime - time.Second)
    if err := isCodeExpired(past); err == nil {
        t.Fatalf("expected expired error, got nil")
    }

    recent := time.Now().Add(-VerificationCodeLifetime + time.Second)
    if err := isCodeExpired(recent); err != nil {
        t.Fatalf("unexpected error for non-expired code: %v", err)
    }
}

func TestValidateCodeRepositoryLayer(t *testing.T) {
    ok := DBCodeData{Code: "123456", CreatedAt: time.Now()}
    if valid, err := validateCodeRepositoryLayer("123456", ok); !valid || err != nil {
        t.Fatalf("expected valid code, got valid=%v err=%v", valid, err)
    }

    mismatch := DBCodeData{Code: "654321", CreatedAt: time.Now()}
    if valid, err := validateCodeRepositoryLayer("123456", mismatch); valid || err == nil {
        t.Fatalf("expected mismatch error, got valid=%v err=%v", valid, err)
    }

    expired := DBCodeData{Code: "123456", CreatedAt: time.Now().Add(-VerificationCodeLifetime - time.Second)}
    if valid, err := validateCodeRepositoryLayer("123456", expired); valid || err == nil {
        t.Fatalf("expected expired error, got valid=%v err=%v", valid, err)
    }
}


