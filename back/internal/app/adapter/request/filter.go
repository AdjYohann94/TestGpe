package request

import (
	"time"
)

type CommonFilters struct {
	ID              []uint     `form:"id[]" filter:"in" field:"id"`
	CreatedAtBefore *time.Time `form:"createdAt[before]" filter:"<" field:"created_at"`
	CreatedAtAfter  *time.Time `form:"createdAt[after]" filter:">" field:"created_at"`
	UpdatedAtBefore *time.Time `form:"updatedAt[before]" filter:"<" field:"updated_at"`
	UpdatedAtAfter  *time.Time `form:"updatedAt[after]" filter:">" field:"updated_at"`
}
