// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package request

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/harness/gitness/internal/api/usererror"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"

	"github.com/go-chi/chi"
)

const (
	PathParamRemainder = "*"

	QueryParamCreatedBy = "created_by"
	QueryParamSort      = "sort"
	QueryParamOrder     = "order"
	QueryParamQuery     = "query"

	QueryParamState = "state"
	QueryParamKind  = "kind"
	QueryParamType  = "type"

	QueryParamAfter  = "after"
	QueryParamBefore = "before"

	QueryParamPage  = "page"
	QueryParamLimit = "limit"
	PerPageDefault  = 30
	PerPageMax      = 100
)

// GetCookie tries to retrive the cookie from the request or returns false if it doesn't exist.
func GetCookie(r *http.Request, cookieName string) (string, bool) {
	cookie, err := r.Cookie(cookieName)
	if errors.Is(err, http.ErrNoCookie) {
		return "", false
	} else if err != nil {
		// this should never happen - documentation and code only return `nil` or `http.ErrNoCookie`
		panic(fmt.Sprintf("unexpected error from request.Cookie(...) method: %s", err))
	}

	return cookie.Value, true
}

// PathParamOrError tries to retrieve the parameter from the request and
// returns the parameter if it exists and is not empty, otherwise returns an error.
func PathParamOrError(r *http.Request, paramName string) (string, error) {
	val, ok := PathParam(r, paramName)
	if !ok {
		return "", usererror.BadRequestf("Parameter '%s' not found in request path.", paramName)
	}

	return val, nil
}

// PathParamOrEmpty retrieves the path parameter or returns an empty string otherwise.
func PathParamOrEmpty(r *http.Request, paramName string) string {
	val, ok := PathParam(r, paramName)
	if !ok {
		return ""
	}

	return val
}

// PathParam retrieves the path parameter or returns false if it exists.
func PathParam(r *http.Request, paramName string) (string, bool) {
	val := chi.URLParam(r, paramName)
	if val == "" {
		return "", false
	}

	return val, true
}

// QueryParam returns the parameter if it exists.
func QueryParam(r *http.Request, paramName string) (string, bool) {
	query := r.URL.Query()
	if !query.Has(paramName) {
		return "", false
	}

	return query.Get(paramName), true
}

// QueryParamList returns list of the parameter values if they exist.
func QueryParamList(r *http.Request, paramName string) ([]string, bool) {
	query := r.URL.Query()
	if !query.Has(paramName) {
		return nil, false
	}

	return query[paramName], true
}

// QueryParamOrDefault retrieves the parameter from the query and
// returns the parameter if it exists, otherwise returns the provided default value.
func QueryParamOrDefault(r *http.Request, paramName string, deflt string) string {
	val, ok := QueryParam(r, paramName)
	if !ok {
		return deflt
	}

	return val
}

// QueryParamOrError tries to retrieve the parameter from the query and
// returns the parameter if it exists, otherwise returns an error.
func QueryParamOrError(r *http.Request, paramName string) (string, error) {
	val, ok := QueryParam(r, paramName)
	if !ok {
		return "", usererror.BadRequestf("Parameter '%s' not found in query.", paramName)
	}

	return val, nil
}

// QueryParamAsPositiveInt64 extracts an integer parameter from the request query.
// If the parameter doesn't exist the provided default value is returned.
func QueryParamAsPositiveInt64OrDefault(r *http.Request, paramName string, deflt int64) (int64, error) {
	value, ok := QueryParam(r, paramName)
	if !ok {
		return deflt, nil
	}

	valueInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil || valueInt <= 0 {
		return 0, usererror.BadRequestf("Parameter '%s' must be a positive integer.", paramName)
	}

	return valueInt, nil
}

// QueryParamAsPositiveInt64 extracts an integer parameter from the request query.
// If the parameter doesn't exist an error is returned.
func QueryParamAsPositiveInt64(r *http.Request, paramName string) (int64, error) {
	value, err := QueryParamOrError(r, paramName)
	if err != nil {
		return 0, err
	}

	valueInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil || valueInt <= 0 {
		return 0, usererror.BadRequestf("Parameter '%s' must be a positive integer.", paramName)
	}

	return valueInt, nil
}

// PathParamAsPositiveInt64 extracts an integer parameter from the request path.
func PathParamAsPositiveInt64(r *http.Request, paramName string) (int64, error) {
	rawValue, err := PathParamOrError(r, paramName)
	if err != nil {
		return 0, err
	}

	valueInt, err := strconv.ParseInt(rawValue, 10, 64)
	if err != nil || valueInt <= 0 {
		return 0, usererror.BadRequestf("Parameter '%s' must be a positive integer.", paramName)
	}

	return valueInt, nil
}

// QueryParamAsBoolOrDefault tries to retrieve the parameter from the query and parse it to bool.
func QueryParamAsBoolOrDefault(r *http.Request, paramName string, deflt bool) (bool, error) {
	rawValue, ok := QueryParam(r, paramName)
	if !ok || len(rawValue) == 0 {
		return deflt, nil
	}

	boolValue, err := strconv.ParseBool(rawValue)
	if err != nil {
		return false, usererror.BadRequestf("Parameter '%s' must be a boolean.", paramName)
	}

	return boolValue, nil
}

// GetOptionalRemainderFromPath returns the remainder ("*") from the path or an empty string if it doesn't exist.
func GetOptionalRemainderFromPath(r *http.Request) string {
	return PathParamOrEmpty(r, PathParamRemainder)
}

// GetRemainderFromPath returns the remainder ("*") from the path or an an error if it doesn't exist.
func GetRemainderFromPath(r *http.Request) (string, error) {
	return PathParamOrError(r, PathParamRemainder)
}

// ParseQuery extracts the query parameter from the url.
func ParseQuery(r *http.Request) string {
	return r.URL.Query().Get(QueryParamQuery)
}

// ParsePage extracts the page parameter from the url.
func ParsePage(r *http.Request) int {
	s := r.URL.Query().Get(QueryParamPage)
	i, _ := strconv.Atoi(s)
	if i <= 0 {
		i = 1
	}
	return i
}

// ParseLimit extracts the limit parameter from the url.
func ParseLimit(r *http.Request) int {
	s := r.URL.Query().Get(QueryParamLimit)
	i, _ := strconv.Atoi(s)
	if i <= 0 {
		i = PerPageDefault
	} else if i > PerPageMax {
		i = PerPageMax
	}
	return i
}

// ParseOrder extracts the order parameter from the url.
func ParseOrder(r *http.Request) enum.Order {
	return enum.ParseOrder(
		r.URL.Query().Get(QueryParamOrder),
	)
}

// ParseSort extracts the sort parameter from the url.
func ParseSort(r *http.Request) string {
	return r.URL.Query().Get(QueryParamSort)
}

// ParsePaginationFromRequest parses pagination related info from the url.
func ParsePaginationFromRequest(r *http.Request) types.Pagination {
	return types.Pagination{
		Page: ParsePage(r),
		Size: ParseLimit(r),
	}
}

// ParseListQueryFilterFromRequest parses pagination and query related info from the url.
func ParseListQueryFilterFromRequest(r *http.Request) types.ListQueryFilter {
	return types.ListQueryFilter{
		Query:      ParseQuery(r),
		Pagination: ParsePaginationFromRequest(r),
	}
}
