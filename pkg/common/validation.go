package common

import uuid "github.com/satori/go.uuid"

func IsValidUUID(u string) bool {
	_, err := uuid.FromString(u)
	return err == nil
}

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
