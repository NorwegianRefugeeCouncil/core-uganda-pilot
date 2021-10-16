package documents

import (
	"github.com/nrc-no/core/pkg/pointers"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func Test_validateDocumentID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "validRootKey",
			id:      "/key",
			wantErr: false,
		}, {
			name:    "validSubKey",
			id:      "/a/b",
			wantErr: false,
		}, {
			name:    "validRootExtension",
			id:      "/a.jpg",
			wantErr: false,
		}, {
			name:    "validSubExtension",
			id:      "/a/b.jpg",
			wantErr: false,
		}, {
			name:    "invalidTrailingCharacter",
			id:      "/abc!",
			wantErr: true,
		}, {
			name:    "invalidLeadingCharacter",
			id:      "!abc",
			wantErr: true,
		}, {
			name:    "invalidInnerLeadingCharacter",
			id:      "/!abc",
			wantErr: true,
		}, {
			name:    "invalidSubLeadingCharacter",
			id:      "/abc/!abc",
			wantErr: true,
		}, {
			name:    "invalidSubTrailingCharacter",
			id:      "/abc/abc!/abc",
			wantErr: true,
		}, {
			name:    "leadingSpace",
			id:      " /abc",
			wantErr: true,
		}, {
			name:    "trailingSpace",
			id:      "/abc ",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateDocumentID(tt.id); (err != nil) != tt.wantErr {
				t.Errorf("validateDocumentID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getDocumentDBFilter(t *testing.T) {
	tests := []struct {
		name string
		args DocumentRef
		want interface{}
	}{
		{
			name: "withVersion",
			args: NewDocumentRef("bucket", "key", pointers.Int64(10)),
			want: bson.M{
				"id":              "key",
				"bucketId":        "bucket",
				"resourceVersion": int64(10),
				"isDeleted":       false,
			},
		}, {
			name: "withoutVersion",
			args: NewDocumentRef("bucket", "key", nil),
			want: bson.M{
				"id":              "key",
				"bucketId":        "bucket",
				"isLatestVersion": true,
				"isDeleted":       false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getDocumentDBFilter(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
