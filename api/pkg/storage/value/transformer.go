package value

import "go.mongodb.org/mongo-driver/bson"

type Transformer interface {
	TransformFromStorage(data []byte) (out []byte, stale bool, err error)
	TransformToStorage(data []byte) (out []byte, err error)
}

type BSONTransformer struct {
}

var _ Transformer = &BSONTransformer{}

func (B BSONTransformer) TransformFromStorage(data []byte) (out []byte, stale bool, err error) {

	var doc = bson.M{}
	bson.Unmarshal(data, &doc)

	panic("implement me")
}

func (B BSONTransformer) TransformToStorage(data []byte) (out []byte, err error) {
	panic("implement me")
}
