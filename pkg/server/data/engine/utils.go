package engine

import (
	"bytes"
)

// joinStrings joins a list of strings into a single string with a separator
// it is a zero-dependant version of strings.Join(strings, sep)
func joinStrings(strings []string, separator string) string {
	if len(strings) == 0 {
		return ""
	}
	if len(strings) == 1 {
		return strings[0]
	}
	n := len(separator) * (len(strings) - 1)
	for i := 0; i < len(strings); i++ {
		n += len(strings[i])
	}
	var b = &StringBuilder{}
	b.buf.Grow(n)
	b.WriteString(strings[0])
	for _, s := range strings[1:] {
		b.WriteString(separator)
		b.WriteString(s)
	}
	return b.String()
}

func repeatStrings(str string, n int) []string {
	var ret []string
	for i := 0; i < n; i++ {
		ret = append(ret, str)
	}
	return ret
}

// StringBuilder is a simple string builder
// it is a zero-dependant version of strings.Builder
type StringBuilder struct {
	buf bytes.Buffer
}

// WriteString writes a string to the buffer
func (b *StringBuilder) WriteString(s string) {
	b.buf.WriteString(s)
}

// String returns the string representation of the buffer
func (b *StringBuilder) String() string {
	return b.buf.String()
}

// Reset resets the buffer
func (b *StringBuilder) Reset() {
	b.buf.Reset()
}

// sortStrings uses bubble sort to sort a list of strings
// it is a zero-dependant version of sort.Strings(strings)
func sortStrings(strings []string) {
	for i := 0; i < len(strings); i++ {
		for j := 0; j < len(strings)-1; j++ {
			if strings[j] > strings[j+1] {
				strings[j], strings[j+1] = strings[j+1], strings[j]
			}
		}
	}
}
