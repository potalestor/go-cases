package activity

import (
	"fmt"
	"testing"
)

func TestActivity_String(t *testing.T) {
	tests := []struct {
		a    Activity
		want string
	}{
		{Auth, "Authentication OK"},
		{Auth | AuthUserNotRegistered, "Authentication failed: user is not registered"},
		{Auth | AuthUserLocked, "Authentication failed: user is locked"},
		{Auth | AuthUserOutsider, "Authentication failed: user is outsider"},
		{Auth | AuthOrgLocked, "Authentication failed: user company is locked"},
		{Auth | AuthCertError, "Authentication failed: user certificate error"},
		{Auth | AuthAnotherError, "Authentication failed: another reason"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := tt.a.String()
			if got != tt.want {
				t.Errorf("Activity.String() = %v, want %v", got, tt.want)
			} else {
				fmt.Printf("Activity.String() = %v\n", got)
			}
		})
	}
}
