package types

import "github.com/nrc-no/core/pkg/api/meta"

const group = "core.nrc.no"
const version = "v1"
const folders = "folders"

var FolderGR = meta.GroupResource{
	Group:    group,
	Resource: folders,
}
