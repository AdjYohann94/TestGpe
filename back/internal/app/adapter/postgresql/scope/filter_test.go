package scope

import (
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"
)

func TestCommonFiltersBuildMapCondition(t *testing.T) {
	date := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	type F struct {
		ID              []uint     `filter:"in" field:"id"`
		CreatedAtBefore *time.Time `filter:"<" field:"created_at"`
		CreatedAtAfter  *time.Time `filter:">" field:"created_at"`
		UpdatedAtBefore *time.Time `filter:"<" field:"updated_at"`
		UpdatedAtAfter  *time.Time `filter:">" field:"updated_at"`
	}
	filters := F{
		ID:              []uint{1, 2, 3},
		CreatedAtBefore: &date,
		CreatedAtAfter:  &date,
		UpdatedAtBefore: &date,
		UpdatedAtAfter:  &date,
	}

	td.CmpMap(t, BuildMapCondition(filters), map[string]interface{}{
		"id in ?":        []uint{1, 2, 3},
		"created_at > ?": &date,
		"created_at < ?": &date,
		"updated_at > ?": &date,
		"updated_at < ?": &date,
	}, td.MapEntries{})
}

func TestCommonFiltersBuildMapConditionEmptyFilters(t *testing.T) {
	date := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	type Filters struct {
		ID              []uint     `filter:"in" field:"id"`
		CreatedAtBefore *time.Time `filter:"<" field:"created_at"`
		CreatedAtAfter  *time.Time `filter:">" field:"created_at"`
		UpdatedAtBefore *time.Time `filter:"<" field:"updated_at"`
		UpdatedAtAfter  *time.Time `filter:">" field:"updated_at"`
	}
	filters := Filters{
		ID:             []uint{1, 2, 3},
		CreatedAtAfter: &date,
	}

	td.CmpMap(t, BuildMapCondition(filters), map[string]interface{}{
		"id in ?":        []uint{1, 2, 3},
		"created_at > ?": &date,
	}, td.MapEntries{})
}

func TestCommonFiltersBuildMapConditionEmptyFiltersWithoutPointers(t *testing.T) {
	type Example struct {
		TestZeroPointer       *int
		TestNilPointer        *int
		TestNonZeroPointer    *int
		TestZeroNonPointer    int
		TestNonZeroNonPointer int
	}

	zero := 0
	nonZero := 1

	td.CmpMap(t, BuildMapCondition(Example{
		TestZeroPointer:       &zero,
		TestNilPointer:        nil,
		TestNonZeroPointer:    &nonZero,
		TestZeroNonPointer:    zero,
		TestNonZeroNonPointer: nonZero,
	}), map[string]interface{}{
		"test_non_zero_non_pointer = ?": 1,
		"test_non_zero_pointer = ?":     &nonZero,
		"test_zero_pointer = ?":         &zero,
	}, td.MapEntries{})
}

func TestCommonFiltersBuildMapConditionWithoutTagsPointerStruct(t *testing.T) {
	type F struct {
		Test *int
	}

	v := 1
	td.CmpMap(t, BuildMapCondition(&F{Test: &v}), map[string]interface{}{
		"test = ?": &v,
	}, td.MapEntries{})
}

func TestCommonFiltersBuildMapConditionWithoutTags(t *testing.T) {
	type F struct {
		Test *int
	}

	v := 1
	td.CmpMap(t, BuildMapCondition(F{Test: &v}), map[string]interface{}{
		"test = ?": &v,
	}, td.MapEntries{})
}

func TestCommonFiltersBuildMapConditionWithoutTagsWithoutPointer(t *testing.T) {
	type F struct {
		Test int
	}

	td.CmpMap(t, BuildMapCondition(F{Test: 1}), map[string]interface{}{
		"test = ?": 1,
	}, td.MapEntries{})
}

func TestCommonFiltersBuildMapConditionWithoutTagsWithoutPointerAndStructSlice(t *testing.T) {
	type F struct {
		Test []int
	}

	td.CmpMap(t, BuildMapCondition(F{Test: []int{1, 2, 3}}), map[string]interface{}{
		"test = ?": []int{1, 2, 3},
	}, td.MapEntries{})
}

func TestToSnakeCase(t *testing.T) {
	td.Cmp(t, CamelCaseToSnakeCase("Test"), "test")
	td.Cmp(t, CamelCaseToSnakeCase("TestCamelCase"), "test_camel_case")
	td.Cmp(t, CamelCaseToSnakeCase(""), "")
	td.Cmp(t, CamelCaseToSnakeCase("_"), "_")
	td.Cmp(t, CamelCaseToSnakeCase("is_snake_case"), "is_snake_case")
	td.Cmp(t, CamelCaseToSnakeCase("notcamel"), "notcamel")
}

func TestBuildFilterInlineCondition(t *testing.T) {
	date := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	type F struct {
		ID              []uint     `filter:"in" field:"id"`
		CreatedAtBefore *time.Time `filter:"<" field:"created_at"`
		CreatedAtAfter  *time.Time `filter:">" field:"created_at"`
		UpdatedAtBefore *time.Time `filter:"<" field:"updated_at"`
		UpdatedAtAfter  *time.Time `filter:">" field:"updated_at"`
	}
	filters := F{
		ID:             []uint{1, 2, 3},
		CreatedAtAfter: &date,
	}

	td.CmpStruct(t, BuildAndFilterInlineCondition(filters), Filters{
		Query: "created_at > ? and id in ?",
		Args:  []interface{}{&date, []uint{1, 2, 3}},
	}, td.StructFields{})
}
