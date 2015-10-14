// Copyright 2014 Martini Authors
// Copyright 2014 The Macaron Authors
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

package method

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gopkg.in/macaron.v1"
)

var tests = []struct {
	Method         string
	OverrideMethod string
	ExpectedMethod string
}{
	{"POST", "PUT", "PUT"},
	{"POST", "PATCH", "PATCH"},
	{"POST", "DELETE", "DELETE"},
	{"GET", "GET", "GET"},
	{"HEAD", "HEAD", "HEAD"},
	{"GET", "PUT", "GET"},
	{"HEAD", "DELETE", "HEAD"},
}

func Test_Override(t *testing.T) {
	for _, test := range tests {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(test.Method, "/", nil)
		if err != nil {
			t.Error(err)
		}
		OverrideRequestMethod(r, test.OverrideMethod)
		Override()(w, r)
		if r.Method != test.ExpectedMethod {
			t.Errorf("Expected %s, got %s", test.ExpectedMethod, r.Method)
		}
	}
}

func selectRoute(m *macaron.Macaron, method string, h macaron.Handler) {
	switch method {
	case "GET":
		m.Get("/", h)
	case "PATCH":
		m.Patch("/", h)
	case "POST":
		m.Post("/", h)
	case "PUT":
		m.Put("/", h)
	case "DELETE":
		m.Delete("/", h)
	case "OPTIONS":
		m.Options("/", h)
	case "HEAD":
		m.Head("/", h)
	default:
		panic("bad method")
	}
}

func Test_SelectiveRouter(t *testing.T) {
	for _, test := range tests {
		w := httptest.NewRecorder()
		m := macaron.New()

		done := make(chan bool)
		selectRoute(m, test.ExpectedMethod, func(rq *http.Request) {
			done <- true
		})

		req, err := http.NewRequest(test.Method, "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		OverrideRequestMethod(req, test.OverrideMethod)

		m.Before(Override())
		go m.ServeHTTP(w, req)
		select {
		case <-done:
		case <-time.After(30 * time.Millisecond):
			t.Errorf("Expected router to route to %s, got something else (%v).", test.ExpectedMethod, test)
		}
	}
}

func Test_In(t *testing.T) {
	for _, test := range tests {
		w := httptest.NewRecorder()
		m := macaron.New()
		m.Before(Override())
		m.Use(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != test.ExpectedMethod {
				t.Errorf("Expected %s, got %s", test.ExpectedMethod, r.Method)
			}
		})

		r, err := http.NewRequest(test.Method, "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		OverrideRequestMethod(r, test.OverrideMethod)

		m.ServeHTTP(w, r)
	}

}

func Test_ParamenterOverride(t *testing.T) {
	for _, test := range tests {
		w := httptest.NewRecorder()
		m := macaron.New()
		m.Before(Override())
		m.Use(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != test.ExpectedMethod {
				t.Errorf("Expected %s, got %s", test.ExpectedMethod, r.Method)
			}
		})

		query := "_method=" + test.OverrideMethod
		r, err := http.NewRequest(test.Method, "/?"+query, nil)
		if err != nil {
			t.Fatal(err)
		}

		m.ServeHTTP(w, r)
	}

}
