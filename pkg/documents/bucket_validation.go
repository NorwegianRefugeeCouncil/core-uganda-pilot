package documents

import (
	"github.com/nrc-no/core/pkg/validation"
	"regexp"
)

func validateBucket(bucket *Bucket) validation.ErrorList {
	result := validation.ErrorList{}
	result = append(result, validateBucketName(validation.NewPath("name"), bucket.Name)...)
	return result
}

var bucketNameRegex = regexp.MustCompile("^[a-zA-Z0-9]+([_-][a-zA-Z0-9]+)*$")
var bucketNameMinLength = 3
var bucketNameMaxLength = 128

func validateBucketName(fieldPath *validation.Path, name string) validation.ErrorList {
	result := validation.ErrorList{}

	if len(name) == 0 {
		result = append(result, validation.Required(fieldPath, name))
	}
	if len(name) < bucketNameMinLength {
		result = append(result, validation.TooShort(fieldPath, name, bucketNameMinLength))
	}
	if len(name) > bucketNameMaxLength {
		result = append(result, validation.TooLong(fieldPath, name, bucketNameMaxLength))
	}
	if !bucketNameRegex.MatchString(name) {
		result = append(result, validation.Invalid(fieldPath, name, "invalid character sequence"))
	}

	return result
}
