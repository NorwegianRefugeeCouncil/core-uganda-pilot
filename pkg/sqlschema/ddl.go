package sqlschema

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

func (d1 DDL) Merge(d2 DDLGenerator) DDL {
	d1.Query += d2.DDL().Query
	d1.Args = append(d1.Args, d2.DDL().Args...)
	return d1
}

func (d1 DDL) MergeAll(separator string, d ...DDLGenerator) DDL {
	for i, generator := range d {
		d1 = d1.Merge(generator)
		if i != len(d)-1 {
			d1 = d1.WriteString(separator)
		}
	}
	return d1
}

func writeParamPlaceholders(count int) string {
	var args []string
	for i := 0; i < count; i++ {
		args = append(args, "?")
	}
	return strings.Join(args, ", ")
}
