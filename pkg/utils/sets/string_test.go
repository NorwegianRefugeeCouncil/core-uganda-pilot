package sets

import (
	"reflect"
	"testing"
)

func TestString_Len(t *testing.T) {
	tests := []struct {
		name string
		s    String
		want int
	}{
		{
			name: "empty",
			s:    NewString(),
			want: 0,
		}, {
			name: "notEmpty",
			s:    NewString("A"),
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_List(t *testing.T) {
	tests := []struct {
		name string
		s    String
		want []string
	}{
		{
			name: "shouldSort",
			s:    NewString("b", "a"),
			want: []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.List(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Has(t *testing.T) {
	tests := []struct {
		name string
		s    String
		args string
		want bool
	}{
		{
			name: "contains",
			s:    NewString("a"),
			args: "a",
			want: true,
		}, {
			name: "notContains",
			s:    NewString("a"),
			args: "b",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Has(tt.args); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_HasAll(t *testing.T) {
	tests := []struct {
		name string
		s    String
		args []string
		want bool
	}{
		{
			name: "containsAll",
			s:    NewString("a", "b"),
			args: []string{"a", "b"},
			want: true,
		}, {
			name: "doesNotContainAll",
			s:    NewString("a", "b"),
			args: []string{"a", "b", "c"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.HasAll(tt.args...); got != tt.want {
				t.Errorf("HasAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_HasAny(t *testing.T) {
	tests := []struct {
		name string
		s    String
		args []string
		want bool
	}{
		{
			name: "doesHaveAny",
			s:    NewString("a", "b"),
			args: []string{"b"},
			want: true,
		}, {
			name: "doesNotHaveAny",
			s:    NewString("a", "b"),
			args: []string{"c"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.HasAny(tt.args...); got != tt.want {
				t.Errorf("HasAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Insert(t *testing.T) {
	tests := []struct {
		name string
		s    String
		args []string
		want String
	}{
		{
			name: "insertOne",
			s:    NewString("a"),
			args: []string{"b"},
			want: NewString("a", "b"),
		}, {
			name: "insertMultiple",
			s:    NewString("a"),
			args: []string{"b", "c"},
			want: NewString("a", "b", "c"),
		}, {
			name: "insertDuplicate",
			s:    NewString("a"),
			args: []string{"a", "a"},
			want: NewString("a"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Insert(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Delete(t *testing.T) {
	tests := []struct {
		name string
		s    String
		args []string
		want String
	}{
		{
			name: "delete",
			s:    NewString("a"),
			args: []string{"a"},
			want: NewString(),
		}, {
			name: "deleteNonExisting",
			s:    NewString("a"),
			args: []string{"b"},
			want: NewString("a"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Delete(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Difference(t *testing.T) {
	tests := []struct {
		name string
		s    String
		args String
		want String
	}{
		{
			name: "allDifferent",
			s:    NewString("a"),
			args: NewString("b"),
			want: NewString("a"),
		}, {
			name: "noneDifferent",
			s:    NewString("a", "b"),
			args: NewString("a", "b"),
			want: NewString(),
		}, {
			name: "oneDifferent",
			s:    NewString("a", "b"),
			args: NewString("b", "c"),
			want: NewString("a"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Difference(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Difference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Union(t *testing.T) {
	tests := []struct {
		name string
		s    String
		args String
		want String
	}{
		{
			name: "empty",
			s:    NewString(),
			args: NewString(),
			want: NewString(),
		}, {
			name: "allSame",
			s:    NewString("a"),
			args: NewString("a"),
			want: NewString("a"),
		}, {
			name: "different",
			s:    NewString("a"),
			args: NewString("b"),
			want: NewString("a", "b"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Union(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Intersection(t *testing.T) {
	tests := []struct {
		name string
		s    String
		args String
		want String
	}{
		{
			name: "empty",
			s:    NewString(),
			args: NewString(),
			want: NewString(),
		}, {
			name: "different",
			s:    NewString("a"),
			args: NewString("b"),
			want: NewString(),
		}, {
			name: "common",
			s:    NewString("a", "b"),
			args: NewString("b", "c"),
			want: NewString("b"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Intersection(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_IsSuperset(t *testing.T) {
	tests := []struct {
		name string
		s    String
		args String
		want bool
	}{
		{
			name: "bothEmpty",
			s:    NewString(),
			args: NewString(),
			want: true,
		}, {
			name: "leftEmpty",
			s:    NewString(),
			args: NewString("b"),
			want: false,
		}, {
			name: "rightEmpty",
			s:    NewString("b"),
			args: NewString(),
			want: true,
		}, {
			name: "same",
			s:    NewString("b"),
			args: NewString("b"),
			want: true,
		}, {
			name: "proper",
			s:    NewString("a", "b"),
			args: NewString("b"),
			want: true,
		}, {
			name: "not",
			s:    NewString("a", "b"),
			args: NewString("b", "c"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsSuperset(tt.args); got != tt.want {
				t.Errorf("IsSuperset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Equal(t *testing.T) {
	tests := []struct {
		name string
		s    String
		args String
		want bool
	}{
		{
			name: "empty",
			s:    NewString(),
			args: NewString(),
			want: true,
		}, {
			name: "notEqual",
			s:    NewString("a"),
			args: NewString("b"),
			want: false,
		}, {
			name: "equal",
			s:    NewString("a"),
			args: NewString("a"),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Equal(tt.args); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
