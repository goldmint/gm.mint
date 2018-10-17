package sumus

import (
	"testing"
)

func TestIsValid58(t *testing.T) {

	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{"ok", "qgQdnYdmnhXmA9N7hDHYVTx1BBmCDpeVnpNb5A8mkBt66PDF4", true},
		{"ok", "RRAeE4H6wMcoYyG3Lymi6UY5VyeupXXgxQrnWFXvgrcqbKwwn", true},
		{"ok", "2C1LhVBGsNrYgYo32ebGZLuQsUXtB9MohWP9ohyoe9DgvJEfmg", true},
		{"ok", "k1yMXnDxUAfHDiGHT2xQrgGU9f6rvBtBuVfcSQi9YQwVXAn5P", true},
		{"ok", "28VY5m11HKiiV7q9J12rQHqYGJKfbrLah8KmNBeaAGgZKmtBCu", true},
		{"ok", "2tZWtWnzPSwwQfsdm5x7TWcfsDfvRfH1hGkfsexAnxmRCS4ybn", true},
		{"ok", "Ys1Tjpn2sft5ktbc6rpjbMdyqThEa49nTH4ij5VMouvwJAQG", true},
		{"ok", "2VE3sWZsGF8kypaP7SXam96rTnxbh7GQLwPikFgbZdMYNEwSx2", true},
		{"fail", "RRAeE4H6wMcoYyG3Lymi6UY5VyeupXXgxQrnWFXvgrcqbKwwo", false},
		{"fail", "Qyd7MtJViy8uUzEUb7UW1oqziXSJYUcVi84xtkZHcKicmHEcH", false},
		{"fail", "RRAeE4H6wMcoYyG3Lymi6UY5VyxupXXgxQrnWFXvgrcqbKwwn", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Unpack58(tt.arg)
			got := err == nil
			if got != tt.want {
				t.Errorf("TestIsValid58() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPack58(t *testing.T) {

	tests := []struct {
		name    string
		args    string
		wantErr bool
	}{
		{"ok", "1", false},
		{"ok", "satana", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Pack58([]byte(tt.args))
			if _, err := Unpack58(got); err != nil {
				t.Errorf("TestPack58() = %v", err)
			}
		})
	}
}

func TestDateFromStamp(t *testing.T) {

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
			got := DateFromStamp(tt.stamp)
			if got.Format("2006-Jan-02 15:04:05") != tt.want {
				t.Errorf("DateFromStamp() = %v, expected %v", got, tt.want)
			}
		})
	}
}
