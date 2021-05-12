package mongo

type event struct {
	key       string
	value     []byte
	prevValue []byte
	rev       uint64
	isDeleted bool
	isCreated bool
}
