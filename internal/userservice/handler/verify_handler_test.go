package handler

import "testing"

func TestIsCodeEmpty(t *testing.T) {
    if err := isCodeEmpty(""); err == nil {
        t.Fatalf("expected error for empty code, got nil")
    }
    if err := isCodeEmpty("123456"); err != nil {
        t.Fatalf("unexpected error for non-empty code: %v", err)
    }
}

func TestIsCodeLengthMismatch(t *testing.T) {
    if err := isCodeLengthMismatch("12345"); err == nil {
        t.Fatalf("expected length mismatch error for 5 digits, got nil")
    }
    if err := isCodeLengthMismatch("1234567"); err == nil {
        t.Fatalf("expected length mismatch error for 7 digits, got nil")
    }
    if err := isCodeLengthMismatch("123456"); err != nil {
        t.Fatalf("unexpected error for correct length: %v", err)
    }
}

func TestIsCodeNotDigitable(t *testing.T) {
    if err := isCodeNotDigitable("12a456"); err == nil {
        t.Fatalf("expected not-digitable error, got nil")
    }
    if err := isCodeNotDigitable("12345!"); err == nil {
        t.Fatalf("expected not-digitable error with symbol, got nil")
    }
    // Arabic-Indic digits are valid according to unicode.IsDigit
    if err := isCodeNotDigitable("١٢٣٤٥٦"); err != nil {
        t.Fatalf("unexpected error for unicode digits: %v", err)
    }
    if err := isCodeNotDigitable("123456"); err != nil {
        t.Fatalf("unexpected error for numeric code: %v", err)
    }
}

func TestIsUserIDNegative(t *testing.T) {
    if err := isUserIDNegative(-1); err == nil {
        t.Fatalf("expected error for negative userID, got nil")
    }
    if err := isUserIDNegative(0); err != nil {
        t.Fatalf("unexpected error for non-negative userID: %v", err)
    }
}

func TestValidateConfirmationCode(t *testing.T) {
    tests := []struct{
        name string
        code string
        wantErr bool
    }{
        {name: "empty", code: "", wantErr: true},
        {name: "short", code: "12345", wantErr: true},
        {name: "long", code: "1234567", wantErr: true},
        {name: "non-digit", code: "12a456", wantErr: true},
        {name: "ok", code: "123456", wantErr: false},
    }

    for _, tc := range tests {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            err := validateConfirmationCode(tc.code)
            if tc.wantErr && err == nil {
                t.Fatalf("expected error, got nil")
            }
            if !tc.wantErr && err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
        })
    }
}


