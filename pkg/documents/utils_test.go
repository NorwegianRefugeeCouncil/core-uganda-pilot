package documents

import "testing"

func Test_validateObjectId(t *testing.T) {
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
			if err := validateObjectId(tt.id); (err != nil) != tt.wantErr {
				t.Errorf("validateObjectId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
