package api

import "encoding/json"

var digitTable = [10]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

// Revision represents the revision of a record
// It contains a revision number and a Hash
// A valid format a revision string
// is "revision:<number>-<hash>"
type Revision struct {
	// Num is the revision number
	Num int
	// Hash is the Hash of the record
	Hash string
}

// MarshalJSON implements the json.Marshaler interface
func (r Revision) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (r *Revision) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	rev, err := ParseRevision(str)
	if err != nil {
		return err
	}
	*r = rev
	return nil
}

// EmptyRevision is the zero value of a Revision
var EmptyRevision = Revision{}

// NewRevision creates a new revision with the
// given number and hash
func NewRevision(num int, hash string) Revision {
	return Revision{num, hash}
}

// String returns the string representation of the revision
// Format: "<number>-<hash>"
func (r Revision) String() string {
	if len(r.Hash) == 0 {
		return ""
	}
	var ret string
	digits := getDigits(r.Num)
	for _, d := range digits {
		ret += digitTable[d]
	}
	ret += "-" + r.Hash
	return ret
}

func (r Revision) IsEmpty() bool {
	return r.Num == 0 && len(r.Hash) == 0
}

// ParseRevision parses a revision string
// the format is <Num>-<Hash>
func ParseRevision(str string) (rev Revision, err error) {
	if len(str) == 0 {
		return EmptyRevision, nil
	}
	// The length of the hash is 32 (md5)
	// The length of the separator is 1
	// The length of the number is at least 1
	// so the minimum length is 34
	if len(str) < 34 {
		return rev, NewError(ErrCodeInvalidRevision, "Invalid revision length")
	}
	hash := str[len(str)-32:]
	// the hash can only contain hex characters
	for _, c := range hash {
		if c < '0' || c > 'f' {
			return rev, NewError(ErrCodeInvalidRevision, "Invalid character in hash")
		}
	}
	// the number is from position 0 to position len(str)-33
	numStr := str[:len(str)-33]
	for i, c := range numStr {
		// the number cannot start with a zero
		if i == 0 && c == '0' {
			return rev, NewError(ErrCodeInvalidRevision, "Invalid character in version")
		}
		// the number can only contain digits
		if c < '0' || c > '9' {
			return rev, NewError(ErrCodeInvalidRevision, "Invalid character in version")
		}
	}
	// the number is parsed as an int
	num, err := parseInt(numStr)
	if err != nil {
		return rev, NewError(ErrCodeInvalidRevision, "Invalid version")
	}
	// check that the separator is correct
	if str[len(numStr)] != '-' {
		return rev, NewError(ErrCodeInvalidRevision, "Invalid revision")
	}
	rev.Num = num
	rev.Hash = hash
	return rev, nil
}
