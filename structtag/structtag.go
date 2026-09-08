package structtag

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
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
		// No explicit tag override — derive the segment from the Go field
		// name itself, splitting multi-word identifiers on word boundaries
		// (EnableBasic -> ENABLE_BASIC) so it reads the way every other
		// nesting-derived segment already does. An explicit tag (the branch
		// above) is left exactly as written; only the auto-derived name is
		// split.
		ret = toSnakeCase(field.Name)
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

// toSnakeCase converts a Go exported identifier such as "EnableBasic",
// "HTTPOnly", or "PublicBaseURL" into an underscore-separated, uppercase
// form ("ENABLE_BASIC", "HTTP_ONLY", "PUBLIC_BASE_URL").
//
// Before this existed, updateTag used the bare field name unmodified, so a
// multi-word field only ever produced one run-together uppercase segment
// (EnableBasic -> ENABLEBASIC) instead of the word-separated form every
// consumer actually expects and every other config field with this style of
// name implicitly promises — env vars for a struct nested two levels deep
// worked fine (e.g. HUNT_MEDIA_DIR, since every segment was already a single
// word), but a multi-word field anywhere in the path silently produced a
// name nothing was ever going to set (e.g. API_AUTH_ENABLE_BASIC and
// API_CORS_ALLOWED_ORIGINS both silently ignored in production).
//
// A boundary is inserted before an uppercase rune that follows a lowercase
// rune (wordWord -> word_Word), and before the last uppercase rune of a run
// that's immediately followed by a lowercase rune (ABCd -> AB_Cd) — the
// second rule is what keeps an acronym like "HTTP" or "URL" together as one
// segment instead of splitting on every letter.
func toSnakeCase(s string) string {
	runes := []rune(s)
	var b strings.Builder
	b.Grow(len(runes) + 4)
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			nextIsLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if unicode.IsLower(prev) || (unicode.IsUpper(prev) && nextIsLower) {
				b.WriteRune('_')
			}
		}
		b.WriteRune(r)
	}
	return strings.ToUpper(b.String())
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
