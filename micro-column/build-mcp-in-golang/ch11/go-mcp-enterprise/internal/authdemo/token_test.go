package authdemo

import (
	"context"
	"net/http"
	"testing"
)

func TestVerifyToken(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		wantUser   string
		wantTenant string
		wantErr    bool
	}{
		{
			name:       "valid alice token",
			token:      "demo-alice-token",
			wantUser:   "alice",
			wantTenant: "acme",
		},
		{
			name:    "invalid token",
			token:   "invalid-token",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := VerifyToken(
				context.Background(),
				tt.token,
				&http.Request{},
			)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("VerifyToken() error = %v", err)
			}

			if info.UserID != tt.wantUser {
				t.Errorf("UserID = %q, want %q", info.UserID, tt.wantUser)
			}

			tenantID, _ := info.Extra["tenant_id"].(string)
			if tenantID != tt.wantTenant {
				t.Errorf("tenant_id = %q, want %q", tenantID, tt.wantTenant)
			}
		})
	}
}
