package user

import (
	"testing"
)

func TestRegister(t *testing.T) {
	test := []struct {
		name    string
		req     RegisterRequest
		wantErr bool
	}{
		{
			name: "success",
			req: RegisterRequest{
				UserAccount:   "test123",
				UserPassword:  "password123",
				CheckPassword: "password123",
			},
			wantErr: false,
		},
		{
			name: "empty account",
			req: RegisterRequest{
				UserAccount:   "",
				UserPassword:  "password123",
				CheckPassword: "password123",
			},
			wantErr: true,
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			userId, err := Register(&tt.req)
			if tt.wantErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.wantErr && userId != 1 {
				t.Errorf("expected userId to be 1 but got %d", userId)
			}
		})
	}
}
