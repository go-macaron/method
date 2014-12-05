// Copyright 2014 martini-contrib/method Authors
// Copyright 2014 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package method is a middleware that implements HTTP method override for Macaron.
package method

import (
	"errors"
	"net/http"

	"github.com/Unknwon/macaron"
)

const (
	// HeaderHTTPMethodOverride is a commonly used
	// Http header to override the method.
	HeaderHTTPMethodOverride = "X-HTTP-Method-Override"
	// ParamHTTPMethodOverride is a commonly used
	// HTML form parameter to override the method.
	ParamHTTPMethodOverride = "_method"
)

var httpMethods = []string{"PUT", "PATCH", "DELETE"}

// ErrInvalidOverrideMethod is returned when
// an invalid http method was given to OverrideRequestMethod.
var ErrInvalidOverrideMethod = errors.New("invalid override method")

func isValidOverrideMethod(method string) bool {
	for _, m := range httpMethods {
		if m == method {
			return true
		}
	}
	return false
}

// OverrideRequestMethod overrides the http
// request's method with the specified method.
func OverrideRequestMethod(r *http.Request, method string) error {
	if !isValidOverrideMethod(method) {
		return ErrInvalidOverrideMethod
	}
	r.Header.Set(HeaderHTTPMethodOverride, method)
	return nil
}

// Override checks for the X-HTTP-Method-Override header
// or the HTML for parameter, `_method`
// and uses (if valid) the http method instead of
// Request.Method.
// This is especially useful for http clients
// that don't support many http verbs.
// It isn't secure to override e.g a GET to a POST,
// so only Request.Method which are POSTs are considered.
func Override() macaron.BeforeHandler {
	return func(w http.ResponseWriter, r *http.Request) bool {
		if r.Method == "POST" {
			m := r.FormValue(ParamHTTPMethodOverride)
			if isValidOverrideMethod(m) {
				OverrideRequestMethod(r, m)
			}
			m = r.Header.Get(HeaderHTTPMethodOverride)
			if isValidOverrideMethod(m) {
				r.Method = m
			}
		}
		return false
	}
}
