package scope

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"sort"
	"strings"
)

type Filters struct {
	Query string
	Args  []interface{}
}

// AutoFilterScore is the automatic scope for input filters.
// The input struct contains query strings and ordered arguments,
// the Where gorm clause is used with these values.
// The query string contains placeholders that are injected by gorm in the escape
// sequence to avoid sql injection.
func AutoFilterScore(filters Filters) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(filters.Query, filters.Args...)
	}
}

// BuildAndFilterInlineCondition this function takes a struct in parameter to generate an inline sql condition
// to filter request from input struct values. This function calls BuildMapCondition functions and then generate
// the Query string and append all arguments in the Filters struct in output.
//
// Each filter is a AND condition with other filters.
func BuildAndFilterInlineCondition(f interface{}) Filters {
	filters := Filters{}
	mapConditions := BuildMapCondition(f)
	keys := make([]string, 0, len(mapConditions))
	for k := range mapConditions {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if filters.Query != "" {
			filters.Query += " and "
		}
		filters.Query += k
		filters.Args = append(filters.Args, mapConditions[k])
	}

	return filters
}

// BuildMapCondition build a mapping interface to provide to gorm a
// valid structure to filter the requested resource.
// This function takes a struct as input and perform analyse on tag to create
// the key filter and value is the struct value.
//
// The field tag is used to generate the corresponding column (if not provided,
// it was the column name lower cased with underscore between caps)
//
// The field filter is used to generate the condition. It was the corresponding
// sql tag (<, >, =, in, >=, <=, like, <>). If not provided, = is used.
//
// Example :
//
// type Scope struct {
//	  Id []uint `filter:"in" field:"id"`
// }
// returns : "id in ?": values
//
// If struct field is not a pointer, nil values for this type are ignored in the output.
// For pointers, they are ignored until they are nil, when the value is not valid (could be zero),
// she's added to the output map.
//
// Warning: if (field, filter) tag tuple is used multiple times, only the last one (struct order) is kept.
// Warning: embedded structs are not supported (dive reading).
func BuildMapCondition(f interface{}) map[string]interface{} {
	mapping := make(map[string]interface{})

	// Remove the pointer if input parameter is an interface pointer.
	v := reflect.ValueOf(f)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		// Skip if pointer is nil (zero values are not skipped)
		if v.Field(i).Type().Kind() == reflect.Ptr {
			if v.Field(i).IsNil() {
				continue
			}
		} else {
			// Skip is the value if zero (non pointers only)
			if v.Field(i).IsZero() {
				continue
			}
		}

		// Take the struct field name if field tag not filled
		field, ok := t.Field(i).Tag.Lookup("field")
		if !ok {
			field = CamelCaseToSnakeCase(t.Field(i).Name)
		}
		// Use the = condition is filter tag not filled
		filter, ok := t.Field(i).Tag.Lookup("filter")
		if !ok {
			filter = "="
		}

		val := v.Field(i).Interface()

		mapping[fmt.Sprintf("%s %s ?", field, filter)] = val
	}

	return mapping
}

// CamelCaseToSnakeCase transforms input string to snake_case syntax.
func CamelCaseToSnakeCase(s string) string {
	s = strings.TrimSpace(s)
	buf := make([]rune, 0, len(s))

	for k, next := range s {
		if k != 0 && isUpper(next) {
			buf = append(buf, '_')
		}
		buf = append(buf, next)
	}

	return strings.ToLower(string(buf))
}

func isUpper(ch rune) bool {
	return ch >= 'A' && ch <= 'Z'
}
