package service

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

type Kin struct {
	R routers.Router
}

// NewKinValidator service performs input request validation from OpenApi V3 documentation.
// The documentation could be found in api/swagger.json directory.
//
// Kin service loads and validate the openapi file and boostrap internal gorillamux router.
func NewKinValidator() Kin {
	doc, err := openapi3.NewLoader().LoadFromFile("api/swagger.json")
	if err != nil {
		panic(fmt.Sprintf("cannot read swagger json : %s", err.Error()))
	}

	ctx := context.Background()
	err = doc.Validate(ctx)
	if err != nil {
		panic(fmt.Sprintf("cannot validate swagger json : %s", err.Error()))
	}

	router, _ := gorillamux.NewRouter(doc)

	return Kin{R: router}
}
