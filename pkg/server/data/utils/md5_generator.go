package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"

	"github.com/nrc-no/core/pkg/server/data/api"
)

// Md5RevGenerator implements the RevisionGenerator api.
// It creates a md5 hash out of a given data set with a number prefix
// It produces a valid Revision string
type Md5RevGenerator struct {
}

// Generate implements the RevisionGenerator.Generate
func (r Md5RevGenerator) Generate(num int, data map[string]interface{}) api.Revision {
	h := md5.New()
	var sortedFields []string
	for key := range data {
		sortedFields = append(sortedFields, key)
	}
	sort.Strings(sortedFields)
	for i, sortedField := range sortedFields {
		fieldValue, _ := data[sortedField]
		h.Write([]byte(sortedField + ":" + strconv.Itoa(i) + ":" + fmt.Sprintf("%v", fieldValue)))
	}
	return api.NewRevision(num, hex.EncodeToString(h.Sum(nil)))
}
