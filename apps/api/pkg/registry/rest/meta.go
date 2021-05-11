package rest

import (
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
)

func FillObjectMetaSystemFields(meta metav1.Object) {
	meta.SetCreationTimestamp(metav1.Now())
	// meta.SetUID(uuid.NewV4().String())
}
