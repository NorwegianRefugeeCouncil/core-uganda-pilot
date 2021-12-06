package slices

import (
	"reflect"
	"testing"
)

func TestReversedStrings(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "length-0",
			args: []string{},
			want: []string{},
		}, {
			name: "length-1",
			args: []string{"a"},
			want: []string{"a"},
		},
		{
			name: "length-2",
			args: []string{"a", "b"},
			want: []string{"b", "a"},
		},
		{
			name: "length-3",
			args: []string{"a", "b", "c"},
			want: []string{"c", "b", "a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReversedStrings(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReversedStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}
