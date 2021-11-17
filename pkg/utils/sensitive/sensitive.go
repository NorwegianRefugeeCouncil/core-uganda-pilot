package sensitive

type Sensitive string

func (s Sensitive) MarshalJSON() ([]byte, error) {
	if len(s) == 0 {
		return []byte(""), nil
	}
	if len(s) > 20 {
		return []byte(s[:5] + "*******"), nil
	}
	return []byte("*****"), nil
}
