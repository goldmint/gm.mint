package mint

import (
	"testing"
)

func Test_maskString6p4(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{"", "1", "1"},
		{"", "abcdef", "abcdef"},
		{"", "abcdef1234", "abcdef1234"},
		{"", "abcdef12345", "abcdef***2345"},
		{"", "YeAHCqTJk4aFnHXGV4zaaf3dTqJkdjQzg8TJENmP3zxDMpa97", "YeAHCq***pa97"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskString6P4(tt.str); got != tt.want {
				t.Errorf("MaskString6P4() = %v, want %v", got, tt.want)
			}
		})
	}
}
