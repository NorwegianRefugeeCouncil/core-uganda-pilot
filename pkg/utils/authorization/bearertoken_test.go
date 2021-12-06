package authorization

import (
	"net/http"
	"testing"
)

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		args    http.Header
		want    string
		wantErr bool
	}{
		{
			name: "valid",
			want: "myToken",
			args: map[string][]string{
				"Authorization": {"Bearer myToken"},
			},
			wantErr: false,
		}, {
			name:    "missing Bearer prefix",
			wantErr: true,
			args: map[string][]string{
				"Authorization": {"myToken"},
			},
		}, {
			name:    "missing token",
			wantErr: true,
			args: map[string][]string{
				"Authorization": {"Bearer"},
			},
		}, {
			name:    "empty",
			wantErr: true,
			args: map[string][]string{
				"Authorization": {""},
			},
		}, {
			name:    "missing",
			wantErr: true,
			args:    map[string][]string{},
		}, {
			name:    "nil",
			wantErr: true,
			args:    nil,
		}, {
			name:    "duplicate",
			wantErr: true,
			args: map[string][]string{
				"Authorization": {"Bearer myToken", "Bearer MyOtherToken"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Header: tt.args,
			}
			got, err := ExtractBearerToken(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractBearerToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
