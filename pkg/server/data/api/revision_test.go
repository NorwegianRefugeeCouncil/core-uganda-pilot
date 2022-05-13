package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseRevision(t *testing.T) {
	tests := []struct {
		name    string
		rev     string
		want    Revision
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid revision",
			rev:  "1-96fc52d8fbf5d2adc6d139cb5b2ea099",
			want: Revision{
				Num:  1,
				Hash: "96fc52d8fbf5d2adc6d139cb5b2ea099",
			},
		}, {
			name: "another valid revision",
			rev:  "1232-96fc52d8fbf5d2adc6d139cb5b2ea099",
			want: Revision{
				Num:  1232,
				Hash: "96fc52d8fbf5d2adc6d139cb5b2ea099",
			},
		}, {
			name:    "extra chars",
			rev:     "1-96fc52d8fbf5d2adc6d139cb5b2ea099-",
			wantErr: assert.Error,
		}, {
			name:    "missing hash",
			rev:     "1-",
			wantErr: assert.Error,
		}, {
			name:    "invalid hash chars",
			rev:     "1-9z9z9z9z9z9z9z9z9z9z9z9z9z9z9z9z",
			wantErr: assert.Error,
		}, {
			name:    "leading 0",
			rev:     "0123-96fc52d8fbf5d2adc6d139cb5b2ea099",
			wantErr: assert.Error,
		}, {
			name:    "invalid hash length",
			rev:     "1-9f2",
			wantErr: assert.Error,
		}, {
			name:    "invalid version number",
			rev:     "abc-9f2",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		if tt.wantErr == nil {
			tt.wantErr = assert.NoError
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRevision(tt.rev)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseRevision(%v)", tt.rev)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ParseRevision(%v)", tt.rev)
		})
	}
}

func Test_revision_String(t *testing.T) {
	type fields struct {
		num  int
		hash string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "revision1",
			fields: fields{
				num:  1,
				hash: "96fc52d8fbf5d2adc6d139cb5b2ea099",
			},
			want: "1-96fc52d8fbf5d2adc6d139cb5b2ea099",
		}, {
			name: "revision2",
			fields: fields{
				num:  1232,
				hash: "96fc52d8fbf5d2adc6d139cb5b2ea099",
			},
			want: "1232-96fc52d8fbf5d2adc6d139cb5b2ea099",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Revision{
				Num:  tt.fields.num,
				Hash: tt.fields.hash,
			}
			assert.Equalf(t, tt.want, r.String(), "String()")
		})
	}
}
