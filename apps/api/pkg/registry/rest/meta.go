package rest

import (
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/satori/go.uuid"
)

func FillObjectMetaSystemFields(meta metav1.Object) {
	meta.SetCreationTimestamp(metav1.Now())
	meta.SetUID(uuid.NewV4().String())
}
