package ec

import (
	"testing"
)

func TestGetTds(t *testing.T) {
	for _, tt := range []struct {
		name    string
		wantErr string
	}{
		{
			name: "pass",
		},
	} {

		tds := New(1, 5.0, 1023)

		t.Run(tt.name, func(t *testing.T) {
			_, err := tds.GetTds(77.0)
			if err != nil {
				t.Errorf("\n%v\n !=\n%v", err, tt.wantErr)
			}
		})
	}
}
