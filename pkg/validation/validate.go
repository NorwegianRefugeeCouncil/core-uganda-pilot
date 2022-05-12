package validation

import "reflect"

func Validate(s any, isNew bool) ErrorList {
	var result ErrorList

	method := reflect.ValueOf(&s).MethodByName("Validate")
	if !method.IsValid() {
		method = reflect.ValueOf(&s).Elem().MethodByName("Validate")
	}
	if method.IsValid() {
		for _, errs := range method.Call([]reflect.Value{reflect.ValueOf(isNew)}) {
			e := errs.Interface().(Error)
			result = append(result, &e)
		}
	}

	fields := reflect.TypeOf(s)
	values := reflect.ValueOf(s)

	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		value := values.Field(i)
		if field.Type.Kind() == reflect.Struct {
			result = append(result, Validate(value.Interface(), isNew)...)
		}
	}

	return result
}
