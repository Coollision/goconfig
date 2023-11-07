package structtag

import (
	"errors"
	"fmt"
	"reflect"
)

// ReflectFunc type used to create functions to parse struct and tags
type ReflectFunc func(
	field *reflect.StructField,
	value *reflect.Value,
	tag string) (err error)

var (
	// ErrNotAPointer error when not a pointer
	ErrNotAPointer = errors.New("not a pointer")

	// ErrNotAStruct error when not a struct
	ErrNotAStruct = errors.New("not a struct")

	// ErrTypeNotSupported error when type not supported
	ErrTypeNotSupported = errors.New("type not supported")

	// ErrUndefinedTag error when Tag var is not defined
	ErrUndefinedTag = errors.New("undefined tag")

	// Tag set the main tag
	Tag string

	// TagDefault set tag default
	TagDefault string

	// TagHelper set tag usage
	TagHelper string

	// TagDisabled used to not process an input
	TagDisabled string

	// TagSeparator separe names on environment variables
	TagSeparator string

	// Prefix is a string that would be placed at the beginning of the generated tags.
	Prefix string

	// ParseMap points to each of the supported types
	ParseMap map[reflect.Kind]ReflectFunc
)

// Setup maps and variables
func Setup() {
	TagDisabled = "-"
	TagSeparator = "_"

	ParseMap = make(map[reflect.Kind]ReflectFunc)

	ParseMap[reflect.Struct] = ReflectStruct
	ParseMap[reflect.Array] = ReflectArray
	ParseMap[reflect.Slice] = ReflectArray
}

// Reset maps caling setup function
func Reset() {
	Setup()
}

// Parse tags on struct instance
func Parse(s interface{}, superTag string) (err error) {
	if Tag == "" {
		err = ErrUndefinedTag
		return
	}

	st := reflect.TypeOf(s)
	if st.Kind() != reflect.Ptr {
		err = ErrNotAPointer
		return
	}

	refField := st.Elem()
	if refField.Kind() != reflect.Struct {
		err = ErrNotAStruct
		return
	}

	refValue := reflect.ValueOf(s).Elem()
	for i := 0; i < refField.NumField(); i++ {
		field := refField.Field(i)
		value := refValue.Field(i)
		kind := field.Type.Kind()

		if field.PkgPath != "" {
			continue
		}

		t := updateTag(&field, superTag)
		if t == "" {
			continue
		}

		f, ok := ParseMap[kind]
		if !ok {
			err = ErrTypeNotSupported
			return
		}

		err = f(&field, &value, t)
		if err != nil {
			return
		}
	}
	return
}

// SetBoolDefaults populates the boolean fields of 's' with cfgDefault values
func SetBoolDefaults(s interface{}, superTag string) (err error) {
	if Tag == "" {
		err = ErrUndefinedTag
		return
	}

	st := reflect.TypeOf(s)
	if st.Kind() != reflect.Ptr {
		err = ErrNotAPointer
		return
	}

	refField := st.Elem()
	if refField.Kind() != reflect.Struct {
		err = ErrNotAStruct
		return
	}

	refValue := reflect.ValueOf(s).Elem()
	for i := 0; i < refField.NumField(); i++ {
		field := refField.Field(i)
		kind := field.Type.Kind()
		value := refValue.Field(i)

		if kind == reflect.Bool {

			if field.PkgPath != "" {
				continue
			}

			t := updateTag(&field, superTag)
			if t == "" {
				continue
			}

			defaultValue := field.Tag.Get(TagDefault)
			v := defaultValue == "true" || defaultValue == "t"
			value.SetBool(v)
		} else if kind == reflect.Struct {
			t := updateTag(&field, superTag)
			if t != "" {
				err := SetBoolDefaults(value.Addr().Interface(), "")
				if err != nil {
					return err
				}
			}
		}
	}
	return
}

func updateTag(field *reflect.StructField, superTag string) (ret string) {
	ret = field.Tag.Get(Tag)
	if ret == TagDisabled {
		ret = ""
		return
	}
	if ret == "" {
		ret = field.Name
	}
	if superTag != "" {
		ret = superTag + TagSeparator + ret
		return
	}
	if Prefix != "" {
		ret = Prefix + TagSeparator + ret
	}
	return
}

// ReflectStruct is called when the Parse encounters a sub-structure in the current structure and then calls Parser again to treat the fields of the sub-structure.
func ReflectStruct(field *reflect.StructField, value *reflect.Value, tag string) (err error) {
	err = Parse(value.Addr().Interface(), tag)
	return
}

// ReflectArray is called when the Parse encounters a sub-array in the current structure and then calls Parser again to treat the fields of the sub-array.
func ReflectArray(field *reflect.StructField, value *reflect.Value, tag string) (err error) {
	req := field.Tag.Get("cfgRequired")
	if req == "true" && value.Len() == 0 {
		err = fmt.Errorf("-%v is required", tag)
		return
	}
	switch value.Type().Elem().Kind() {
	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Ptr, reflect.Interface:
		for i := 0; i < value.Len(); i++ {
			err = Parse(value.Index(i).Addr().Interface(), fmt.Sprintf("%s[%d]", tag, i))
			if err != nil {
				return
			}
		}
	}
	return
}
