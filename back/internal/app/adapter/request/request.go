package request

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommonRequest struct{ *gin.Context }

type ClaimsHeaders struct {
	UserID uint
}

type Identifier struct {
	ID uint `uri:"id"`
}

// GetValidatedClaimsHeaders returns the claims from accessToken and refresh token.
// This information was linked to the user connected.
func (r CommonRequest) GetValidatedClaimsHeaders() (*ClaimsHeaders, error) {
	var req ClaimsHeaders
	userID := r.GetString("userID")
	if userID == "" {
		return nil, fmt.Errorf("userID not defined")
	}
	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("userID invalid")
	}
	req.UserID = uint(id)

	return &req, nil
}

// GetValidatedResourceIdentifier returns the resource id from the request url
func (r CommonRequest) GetValidatedResourceIdentifier() (*Identifier, error) {
	var req Identifier
	if err := r.ShouldBindUri(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

// GetValidatedCommonFilters returns the standard field filters
func (r CommonRequest) GetValidatedCommonFilters() (*CommonFilters, error) {
	var req CommonFilters
	if err := r.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	return &req, nil
}
