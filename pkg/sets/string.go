package sets

import (
	"reflect"
	"sort"
	"strings"
)

type Empty struct{}

type String map[string]Empty

type stringSet struct {
	vals map[string]struct{}
}

func NewString(items ...string) String {
	ss := String{}
	ss.Insert(items...)
	return ss
}

func StringKeySet(theMap interface{}) String {
	v := reflect.ValueOf(theMap)
	ret := String{}
	for _, keyValue := range v.MapKeys() {
		ret.Insert(keyValue.Interface().(string))
	}
	return ret
}

func (s String) Len() int {
	return len(s)
}

func (s String) List() []string {
	result := make([]string, len(s), len(s))
	i := 0
	for key := range s {
		result[i] = key
		i++
	}
	sort.Strings(result)
	return result
}

func (s String) ListIntf() []interface{} {
	strList := s.List()
	result := make([]interface{}, len(strList), len(strList))
	for i, value := range strList {
		result[i] = value
	}
	return result
}

func (s String) Join(separator string) string {
	return strings.Join(s.List(), separator)
}

func (s String) UnsortedList() []string {
	result := make([]string, len(s), len(s))
	i := 0
	for key := range s {
		result[i] = key
		i++
	}
	return result
}

func (s String) Has(v string) bool {
	if _, ok := s[v]; ok {
		return true
	}
	return false
}

func (s String) HasAll(items ...string) bool {
	for _, item := range items {
		if _, ok := s[item]; !ok {
			return false
		}
	}
	return true
}

func (s String) HasAny(items ...string) bool {
	for _, item := range items {
		if _, ok := s[item]; ok {
			return true
		}
	}
	return false
}

func (s String) Insert(items ...string) String {
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

func (s String) Delete(items ...string) String {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

func (s String) Difference(s2 String) String {
	result := NewString()
	for key := range s {
		if !s2.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

func (s String) Union(s2 String) String {
	result := NewString()
	for key := range s {
		result.Insert(key)
	}
	for key := range s2 {
		result.Insert(key)
	}
	return result
}

func (s String) Intersection(s2 String) String {
	var walk, other String
	result := NewString()
	if s.Len() < s2.Len() {
		walk = s
		other = s2
	} else {
		walk = s2
		other = s
	}
	for key := range walk {
		if other.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

func (s String) IsSuperset(s2 String) bool {
	for key := range s2 {
		if !s.Has(key) {
			return false
		}
	}
	return true
}

func (s String) Equal(s2 String) bool {
	return s.Len() == s2.Len() && s.IsSuperset(s2)
}

func (s String) PopAny() (string, bool) {
	for key := range s {
		s.Delete(key)
		return key, true
	}
	var zeroValue string
	return zeroValue, false
}
