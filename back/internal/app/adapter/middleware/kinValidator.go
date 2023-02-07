package middleware

import (
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"net/http"

	"gpe_project/internal/app/adapter/service"

	"github.com/gin-gonic/gin"
)

// Kin middleware performs input request validation from service.Kin configuration.
//
// This middleware checks in first time in the request path and method is defined in openapi
// documentation (by matching mux router). If the request does not match any route, 404 error page
// is sent back to client.
//
// Then, is request go through first step, the input body, query parameters, headers, content type and so on
// are validated against the defined expected types, regular expressions, enums...
// If the input body does is incorrect, http.StatusUnprocessableEntity is sent back to client.
func Kin(kin service.Kin) gin.HandlerFunc {
	return func(c *gin.Context) {
		route, params, err := kin.R.FindRoute(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "route not found"})
			return
		}
		if err := openapi3filter.ValidateRequest(c, &openapi3filter.RequestValidationInput{
			Request:    c.Request,
			PathParams: params,
			Route:      route,
		}); err != nil {
			var schemaError *openapi3.SchemaError
			if errors.As(err, &schemaError) {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": schemaError.Reason})
			} else {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			}
			return
		}

		c.Next()
	}
}
