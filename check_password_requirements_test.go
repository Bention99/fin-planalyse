package main

import "testing"

func TestCheckPasswordRequirements(t *testing.T) {
	tests := []struct {
		name    string
		pw      string
		wantErr bool
	}{
		{
			name:    "too short",
			pw:      "Ab1!",
			wantErr: true,
		},
		{
			name:    "exactly 8 chars but missing special",
			pw:      "Abcdef12",
			wantErr: true,
		},
		{
			name:    "missing letter",
			pw:      "1234567!",
			wantErr: true,
		},
		{
			name:    "missing number",
			pw:      "Password!",
			wantErr: true,
		},
		{
			name:    "missing special character",
			pw:      "Password1",
			wantErr: true,
		},
		{
			name:    "valid simple",
			pw:      "Password1!",
			wantErr: false,
		},
		{
			name:    "valid with multiple specials",
			pw:      "a1!aaaaa!",
			wantErr: false,
		},
		{
			name:    "valid with unicode letter",
			pw:      "äbcdef1!",
			wantErr: false,
		},
		{
			name:    "only special and digits no letters",
			pw:      "1234!!!!",
			wantErr: true,
		},
		{
			name:    "only letters and special no digits",
			pw:      "abcdefgh!",
			wantErr: true,
		},
		{
			name:    "only letters and digits no special",
			pw:      "abcd1234",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := checkPasswordRequirements(tc.pw)
			gotErr := err != nil
			if gotErr != tc.wantErr {
				t.Fatalf("checkPasswordRequirements(%q) err=%v, wantErr=%v", tc.pw, err, tc.wantErr)
			}
		})
	}
}