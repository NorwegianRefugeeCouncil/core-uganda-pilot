package fields

type Selector interface{}

type AndSelector []Selector

type OrSelector []Selector

type EqualSelector struct {
	Key   string
	Value interface{}
}

type Fields interface {
}

func And(selectors ...Selector) AndSelector {
	return selectors
}
func Or(selectors ...Selector) OrSelector {
	return selectors
}
func Equal(key string, value interface{}) EqualSelector {
	return EqualSelector{
		Key:   key,
		Value: value,
	}
}
