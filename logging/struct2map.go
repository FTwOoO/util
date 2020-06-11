package logging

import (
	"reflect"
	"strings"
)

func structToMap(obj interface{}) map[string]interface{} {

	var obj1 reflect.Type
	var obj2 reflect.Value
	if reflect.ValueOf(obj).Type().Kind() == reflect.Ptr {
		obj2 = reflect.ValueOf(obj).Elem()
		obj1 = obj2.Type()
	} else {
		obj1 = reflect.TypeOf(obj)
		obj2 = reflect.ValueOf(obj)
	}

	var data = make(map[string]interface{})

	for i := 0; i < obj2.NumField(); i++ {
		tag := obj1.Field(i).Tag.Get("json")
		arr := strings.Split(tag, ",")
		fieldName := arr[0]

		//过滤掉空字符串
		if len(arr) == 2 && arr[1] == "omitempty" {
			if obj1.Field(i).Type.Kind() == reflect.String && obj2.Field(i).Interface() == "" {
				continue
			}
		}
		data[fieldName] = obj2.Field(i).Interface()
	}
	return data
}
