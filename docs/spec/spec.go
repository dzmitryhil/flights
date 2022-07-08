// Package spec defines the open API specification.
//
// Documentation
//
// Schemes: http
// BasePath: /
// Version: 1.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package spec

import "github.com/dzmitryhil/flights/handler"

type headerAccessControlAllowOrigin struct {
	// in:header
	F string `json:"Access-Control-Allow-Origin"`
}

type responseGeneric struct {
	headerAccessControlAllowOrigin
}

// swagger:response responseGenericError400
type responseGenericError400 struct {
	responseGeneric
	// in:body
	F string `json:"body"`
}

// swagger:response responseGenericError500
type responseGenericError500 struct {
	responseGeneric
	// in:body
	F string `json:"body"`
}

// swagger:route POST /flights/path flight flightpath
//
// Returns 200 and path is the path is found.
//
// Responses:
// 200: responseFlightpath200
// 400: responseGenericError400
// 500: responseGenericError500

// swagger:parameters flightpath
type parametersFlightPath200 struct {
	// in:body
	F handler.PostFlightRequestBody `json:"body"`
}

// Information about the path.
// swagger:response responseFlightpath200
type responseFlightpath200 struct {
	responseGeneric
	// in:body
	F handler.PostFlightResponseBody `json:"body"`
}
