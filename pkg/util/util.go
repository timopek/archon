/*
Copyright 2016 The Archon Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"fmt"
	"reflect"
	"strconv"
)

// Convert a string map to a struct with type assertions along the way.
func MapToStruct(in map[string]string, out interface{}, keyPrefix string) (err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("Error: %+v", err2)
		}
	}()

	if in == nil {
		return fmt.Errorf("Can't convert to struct. in is nil")
	}

	ptrType := reflect.TypeOf(out)
	if ptrType.Kind() != reflect.Ptr {
		return fmt.Errorf("Can't convert to struct. out has to be a pointer to a struct")
	}

	v := reflect.ValueOf(out).Elem()
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("Can't convert to struct. out has to be a pointer to a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)

		tagName := fieldType.Tag.Get("k8s")
		if tagName == "" {
			continue
		}
		keyName := keyPrefix + fieldType.Tag.Get("k8s")
		valueString, ok := in[keyName]
		if !ok {
			continue
		}

		var value interface{}
		switch fieldType.Type.Kind() {
		case reflect.Int:
			value, err = strconv.Atoi(valueString)
		case reflect.Bool:
			value, err = strconv.ParseBool(valueString)
		case reflect.Float32:
			value, err = strconv.ParseFloat(valueString, 32)
			value = float32(value.(float64))
		case reflect.Float64:
			value, err = strconv.ParseFloat(valueString, 64)
		case reflect.String:
			value = valueString
		default:
			err = fmt.Errorf("Unsupported field type: %+v", fieldType.Type)
		}
		if err != nil {
			return fmt.Errorf("Unmatch type of field %s: %+v", fieldType.Name, err)
		}
		fieldValue.Set(reflect.ValueOf(value))
	}

	return nil
}

func StructToMap(in interface{}, out map[string]string, prefix string) (err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("Error: %+v", err2)
		}
	}()

	if out == nil {
		return fmt.Errorf("Can't convert to nil map. out is nil")
	}

	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("Can't convert struct. in has to be a struct or *struct")
	}

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)

		tagName := fieldType.Tag.Get("k8s")
		if tagName == "" {
			continue
		}

		key := prefix + tagName
		out[key] = fmt.Sprintf("%v", fieldValue.Interface())
	}

	return nil
}

// Copy the values of maps of the same type
func MapCopy(dst interface{}, src interface{}) (err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("Error: %+v", err2)
		}
	}()

	if src == nil {
		return
	}

	if dst == nil && src != nil {
		return fmt.Errorf("Dst is nil but src is not")
	}

	srct := reflect.TypeOf(src)
	dstt := reflect.TypeOf(dst)
	if srct.Kind() != reflect.Map ||
		dstt.Kind() != reflect.Map ||
		srct.Key().Kind() != dstt.Key().Kind() ||
		srct.Elem().Kind() != dstt.Elem().Kind() {
		return fmt.Errorf("dst and src types don't match")
	}

	srcv := reflect.ValueOf(src)
	dstv := reflect.ValueOf(dst)
	for _, k := range srcv.MapKeys() {
		dstv.SetMapIndex(k, srcv.MapIndex(k))
	}

	return
}

// Check if two maps have the same key and values, an empty map == nil
func MapEqual(dst interface{}, src interface{}) (ret bool) {
	defer func() {
		if err2 := recover(); err2 != nil {
			ret = false
		}
	}()

	srct := reflect.TypeOf(src)
	dstt := reflect.TypeOf(dst)
	if srct.Kind() != reflect.Map ||
		dstt.Kind() != reflect.Map ||
		srct.Key().Kind() != dstt.Key().Kind() ||
		srct.Elem().Kind() != dstt.Elem().Kind() {
		return false
	}

	srcv := reflect.ValueOf(src)
	dstv := reflect.ValueOf(dst)
	if len(srcv.MapKeys()) == 0 && len(dstv.MapKeys()) == 0 {
		return true
	}

	return reflect.DeepEqual(src, dst)
}
