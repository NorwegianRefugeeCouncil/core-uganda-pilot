package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Record_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		name      string
		input     []byte
		expected  Record
		expectErr assert.ErrorAssertionFunc
	}
	testCases := []testCase{
		{
			name:  "empty",
			input: []byte(`{}`),
			expected: Record{
				Attributes: Attributes{},
			},
		}, {
			name:  "with id",
			input: []byte(`{"id":"abc"}`),
			expected: Record{
				ID:         "abc",
				Attributes: Attributes{},
			},
		}, {
			name:  "with revision",
			input: []byte(`{"revision": "1-96fc52d8fbf5d2adc6d139cb5b2ea099"}`),
			expected: Record{
				Revision: Revision{
					Num:  1,
					Hash: "96fc52d8fbf5d2adc6d139cb5b2ea099",
				},
				Attributes: Attributes{},
			},
		}, {
			name: "with data",
			input: []byte(`
			{
				"attributes": {
					"foo":{"string":"bar"}
				}
			}`),
			expected: Record{
				Attributes: Attributes{
					"foo": NewStringValue("bar", true),
				},
			},
		}, {
			name: "full",
			input: []byte(`{
				"id": "abc",
				"revision": "1-96fc52d8fbf5d2adc6d139cb5b2ea099",
                "attributes": {
					"foo": {"string":"bar"},
					"bar": {"null":true}
				}
			}`),
			expected: Record{
				ID: "abc",
				Revision: Revision{
					Num:  1,
					Hash: "96fc52d8fbf5d2adc6d139cb5b2ea099",
				},
				Attributes: Attributes{
					"foo": NewStringValue("bar", true),
					"bar": NewNullValue(),
				},
			},
		},
	}
	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			record := Record{}
			err := json.Unmarshal(testCase.input, &record)
			if testCase.expectErr == nil {
				testCase.expectErr = assert.NoError
			}
			if !testCase.expectErr(t, err) {
				return
			}
			if err != nil {
				return
			}
			assert.Equal(t, testCase.expected, record)
		})
	}
}

func Test_Record_MarshalJSON(t *testing.T) {
	type testCase struct {
		name      string
		columns   map[string]ValueKind
		input     Record
		expected  string
		expectErr assert.ErrorAssertionFunc
	}
	testCases := []testCase{
		{
			name:     "empty",
			input:    Record{},
			expected: `{"id":"","table":"","revision":"","attributes":null}`,
		}, {
			name: "with id",
			input: Record{
				ID: "abc",
			},
			expected: `{"id":"abc","table":"","revision":"","attributes":null}`,
		}, {
			name: "with revision",
			input: Record{
				Revision: Revision{
					Num:  1,
					Hash: "96fc52d8fbf5d2adc6d139cb5b2ea099",
				},
			},
			expected: `{"id":"","table":"","revision":"1-96fc52d8fbf5d2adc6d139cb5b2ea099","attributes":null}`,
		}, {
			name: "with table",
			input: Record{
				Table: "abc",
			},
			expected: `{"id":"","table":"abc","revision":"","attributes":null}`,
		}, {
			name: "with fields",
			input: Record{
				Attributes: Attributes{
					"bar": NewStringValue("", false),
					"foo": NewStringValue("bar", true),
				},
			},
			expected: `{"id":"","table":"","revision":"","attributes":{"bar":{"null":true},"foo":{"string":"bar"}}}`,
		}, {
			name: "full",
			input: Record{
				ID: "abc",
				Revision: Revision{
					Num:  1,
					Hash: "96fc52d8fbf5d2adc6d139cb5b2ea099",
				},
				Table: "foo",
				Attributes: Attributes{
					"bar": NewNullValue(),
					"foo": NewStringValue("bar", true),
				},
			},
			expected: `{"id":"abc","table":"foo","revision":"1-96fc52d8fbf5d2adc6d139cb5b2ea099","attributes":{"bar":{"null":true},"foo":{"string":"bar"}}}`,
		},
	}
	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := json.Marshal(testCase.input)
			if testCase.expectErr == nil {
				testCase.expectErr = assert.NoError
			}
			if !testCase.expectErr(t, err) {
				return
			}
			if err != nil {
				return
			}
			assert.Equal(t, testCase.expected, string(actual))
		})
	}
}
