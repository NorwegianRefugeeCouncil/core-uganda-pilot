package api

type PutRecordOptions struct {
	IsNew bool
}

type PutRecordOption func(o *PutRecordOptions)

var IsNew = func(isNew bool) PutRecordOption {
	return func(o *PutRecordOptions) {
		o.IsNew = isNew
	}
}
