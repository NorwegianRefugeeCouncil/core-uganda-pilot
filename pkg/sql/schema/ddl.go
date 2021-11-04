package schema

import (
	"fmt"
	"strconv"
	"strings"
)

type DDLGenerator interface {
	DDL() DDL
}

type DDL struct {
	Query string
	Args  []interface{}
}

func (d DDL) String() string {
	return fmt.Sprintf("Query:\n%s\nArgs:\n%v", d.Query, d.Args)
}

func (d DDL) DDL() DDL {
	return d
}

func (d DDL) Write(s string, args ...interface{}) DDL {

	parts := strings.Split(s, "?")

	newQry := ""
	for i, part := range parts {
		newQry += part
		if i == len(parts)-1 {
			break
		}
		newQry += "$" + strconv.Itoa(len(d.Args)+1+i)
	}

	d.Query += newQry
	d.Args = append(d.Args, args...)
	return d
}

func (d DDL) WriteF(s string, args ...interface{}) DDL {
	d.Query += fmt.Sprintf(s, args...)
	return d
}

func (d DDL) WriteString(s string) DDL {
	d.Query += s
	return d
}

func (d DDL) WritePlaceholders(args ...interface{}) DDL {
	d.Query += writeParamPlaceholders(len(args))
	d.Args = append(d.Args, args...)
	return d
}

func (d DDL) WriteStringPlaceholders(strs ...string) DDL {
	d.Query += writeParamPlaceholders(len(strs))
	var args []interface{}
	for _, str := range strs {
		args = append(args, str)
	}
	d.Args = append(d.Args, args...)
	return d
}

func (d DDL) WriteStringBefore(s string) DDL {
	d.Query = s + d.Query
	return d
}

func (d DDL) WriteBefore(s string, args ...interface{}) DDL {
	d.Query = s + d.Query
	d.Args = append(args, d.Args...)
	return d
}

func (d DDL) WriteBeforeF(s string, args ...interface{}) DDL {
	d.Query = fmt.Sprintf(s, args...) + d.Query
	return d
}

func (d DDL) WriteArgs(args ...interface{}) DDL {
	d.Args = append(d.Args, args...)
	return d
}

func NewDDL(query string, args ...interface{}) DDL {
	return DDL{
		Query: query,
		Args:  args,
	}
}

func EmptyDDL() DDL {
	return DDL{}
}

func (d DDL) Merge(d2 DDLGenerator) DDL {
	d.Query += d2.DDL().Query
	d.Args = append(d.Args, d2.DDL().Args...)
	return d
}

func (d DDL) MergeAll(separator string, generators ...DDLGenerator) DDL {
	ddl := d
	for i, generator := range generators {
		ddl = ddl.Merge(generator)
		if i != len(generators)-1 {
			ddl = ddl.WriteString(separator)
		}
	}
	return ddl
}

func writeParamPlaceholders(count int) string {
	var args []string
	for i := 0; i < count; i++ {
		args = append(args, "?")
	}
	return strings.Join(args, ", ")
}
