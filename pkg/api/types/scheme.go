package types

import "github.com/nrc-no/core/pkg/api/meta"

const group = "core.nrc.no"
const folders = "folders"
const databases = "databases"
const records = "records"
const forms = "forms"

var FolderGR = meta.GroupResource{
	Group:    group,
	Resource: folders,
}

var DatabaseGR = meta.GroupResource{
	Group:    group,
	Resource: databases,
}

var RecordGR = meta.GroupResource{
	Group:    group,
	Resource: records,
}

var FormGR = meta.GroupResource{
	Group:    group,
	Resource: forms,
}
