package authorization

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenAudience_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		want    TokenAudience
		json    string
		wantErr bool
	}{
		{
			name:    "valid string",
			want:    TokenAudience{"abc"},
			json:    "abc",
			wantErr: false,
		}, {
			name:    "valid list",
			want:    TokenAudience{"abc", "def"},
			json:    `["abc", "def"]`,
			wantErr: false,
		}, {
			name:    "empty list",
			want:    TokenAudience{},
			json:    `[]`,
			wantErr: false,
		}, {
			name:    "empty string",
			want:    TokenAudience{},
			json:    ``,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var aud = TokenAudience{}
			err := aud.UnmarshalJSON([]byte(tt.json))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tt.want, aud)
		})
	}
}
