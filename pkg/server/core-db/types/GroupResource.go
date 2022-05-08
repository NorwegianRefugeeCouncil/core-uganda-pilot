package types

import "github.com/nrc-no/core/pkg/api/meta"

const group = "core.nrc.no"
const entity = "entity"

var EntityGR = meta.GroupResource{
	Group:    group,
	Resource: entity,
}
