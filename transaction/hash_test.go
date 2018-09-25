package transaction

import (
	"testing"

	sumus "github.com/void616/gm-sumus-lib"
)

func TestUnpackHash(t *testing.T) {
	tests := []struct {
		name      string
		hash      string
		wantAddr  string
		wantNonce uint64
		wantErr   bool
	}{
		{"ok", "cqG4tLhKKNd4ZirnFv7HqaYKDdD6c8GuUXdoWwgE6TmBZ6eu885fgkT2BEoJ", "qY4dBwxN7LfAjNeVhoJfKsAk8DjtCY9WGBMTeqvRvBJqcThNp", 1, false},
		{"fail", "2XfAbdqgBp69XHZfFPJH54XY4Rh6qPpKXG8e8YK6BgG6yQgBjmdvYJGGZDsrg1BRmjPHq3M7D2H6QsZ3YH2i", "qY4dBwxN7LfAjNeVhoJfKsAk8DjtCY9WGBMTeqvRvBJqcThNp", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddr, gotNonce, gotErr := UnpackHash(tt.hash)

			if (gotErr != nil) != tt.wantErr {
				t.Errorf("UnpackHash() got err %v, want %v", (gotErr != nil), tt.wantErr)
			}
			if gotErr == nil {
				if sumus.Pack58(gotAddr) != tt.wantAddr {
					t.Errorf("UnpackHash() got addr %v, want %v", sumus.Pack58(gotAddr), tt.wantAddr)
				} else if gotNonce != tt.wantNonce {
					t.Errorf("UnpackHash() got nonce %v, want %v", gotNonce, tt.wantNonce)
				}
			}
		})
	}
}
