package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDigits(t *testing.T) {
	tests := []struct {
		name string
		i    int
		want []int
	}{
		{
			name: "0",
			i:    0,
			want: []int{0},
		}, {
			name: "1",
			i:    1,
			want: []int{1},
		}, {
			name: "10",
			i:    10,
			want: []int{1, 0},
		}, {
			name: "105",
			i:    105,
			want: []int{1, 0, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getDigits(tt.i), "getDigits(%v)", tt.i)
		})
	}
}
