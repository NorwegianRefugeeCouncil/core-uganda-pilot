package pointers

func Int(val int) *int {
	return &val
}

func Int64(val int64) *int64 {
	return &val
}

func Int32(val int32) *int32 {
	return &val
}

func String(val string) *string {
	return &val
}
