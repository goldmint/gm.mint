package mint

import (
	"testing"
)

func TestStampToTime(t *testing.T) {

	tests := []struct {
		name  string
		stamp uint64
		want  string
	}{
		{"ok", uint64(19502164800000000), "2017-Dec-31 12:00:00"},
		{"ok", uint64(19527035308000000), "2018-Oct-15 08:28:28"},
		{"ok", uint64(19527035428000000), "2018-Oct-15 08:30:28"},
		{"ok", uint64(19527219262000000), "2018-Oct-17 11:34:22"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StampToTime(tt.stamp)
			if got.Format("2006-Jan-02 15:04:05") != tt.want {
				t.Errorf("DateToStamp() = %v, expected %v", got, tt.want)
			}
		})
	}
}
