package sqlschema

import (
	"fmt"
	"github.com/lib/pq"
)

type SQLSchema struct {
	Name string
}

func (s SQLSchema) DDL() DDL {
	return NewDDL(fmt.Sprintf("CREATE SCHEMA %s;", pq.QuoteIdentifier(s.Name)))
}
