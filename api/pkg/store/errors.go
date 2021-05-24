package store

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func IsNotFound(err error) bool {
	return err.Error() == "mongo: no documents in result"
}

func IsAlreadyExists(err error) bool {
	return err.Error() == "key already exists"
}

func InterpretGetError(err error, qualifiedResource schema.GroupResource, name string) error {
	switch {
	case IsNotFound(err):
		return errors.NewNotFound(qualifiedResource, name)
	}
	return err
}

func InterpretDeleteError(err error, qualifiedResource schema.GroupResource, name string) error {
	switch {
	case IsNotFound(err):
		return errors.NewNotFound(qualifiedResource, name)
	}
	return err
}

func InterpretUpdateError(err error, qualifiedResource schema.GroupResource, name string) error {
	switch {
	case IsNotFound(err):
		return errors.NewNotFound(qualifiedResource, name)
	}
	return err
}

func InterpretCreateError(err error, qualifiedResource schema.GroupResource, name string) error {
	switch {
	case IsAlreadyExists(err):
		return errors.NewAlreadyExists(qualifiedResource, name)
	default:
		return err
	}
}

func InterpretListError(err error, qualifiedResource schema.GroupResource) error {
	return err
}
