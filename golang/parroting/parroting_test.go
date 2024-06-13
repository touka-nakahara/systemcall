package parroting

import "testing"

func TestParroting(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "hoge"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Parroting()
		})
	}
}
