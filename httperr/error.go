// SPDX-License-Identifier: MIT
//
// Copyright 2025 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

// Package httperr provides the ability to attach and extract the status code from HTTP errors.
package httperr

import (
	"errors"
	"net/http"

	"bursavich.dev/errcode"
	"google.golang.org/grpc/codes"
)

// An Error is an error with an explicit HTTP code.
type Error interface {
	HTTPCode() int
	error
}

// New wraps the given error and adds an HTTP code.
func New(code int, err error) error {
	return &codedError{code, err}
}

type codedError struct {
	code int
	err  error
}

func (ce *codedError) HTTPCode() int { return ce.code }
func (ce *codedError) Error() string { return ce.err.Error() }
func (ce *codedError) Unwrap() error { return ce.err }

var errorCoder errcode.ErrorCoder = errcode.FromFunc(ErrorCode)

// ErrorCoder return the HTTP ErrorCoder.
func ErrorCoder() errcode.ErrorCoder {
	return errorCoder
}

// ErrorCode returns the gRPC code associated with the given error
// if it implements the httperr.Error interface.
func ErrorCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if e, ok := err.(Error); ok || errors.As(err, &e) {
		return ToGRPC(e.HTTPCode())
	}
	return codes.Unknown
}

// ToGRPC returns the gRPC status code associated with the given HTTP status code.
func ToGRPC(httpCode int) codes.Code {
	if 200 <= httpCode && httpCode <= 299 {
		return codes.OK
	}
	switch httpCode {
	case http.StatusBadRequest: // 400
		return codes.InvalidArgument
	case http.StatusUnauthorized: // 401
		return codes.Unauthenticated
	case http.StatusForbidden: // 403
		return codes.PermissionDenied
	case http.StatusNotFound: // 404
		return codes.NotFound
	case http.StatusConflict:
		return codes.Aborted
	case http.StatusRequestedRangeNotSatisfiable: // 416
		return codes.OutOfRange
	case http.StatusTooManyRequests: // 429
		return codes.ResourceExhausted
	case 499:
		return codes.Canceled
	case http.StatusInternalServerError: // 500
		return codes.Internal
	case http.StatusNotImplemented: // 501
		return codes.Unimplemented
	case http.StatusServiceUnavailable: //503
		return codes.Unavailable
	case http.StatusBadGateway: // 504
		return codes.DeadlineExceeded
	}
	return codes.Unknown
}
