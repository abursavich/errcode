// SPDX-License-Identifier: MIT
//
// Copyright 2025 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

// Package grpcerr provides the ability to extract the status code from gRPC errors.
package grpcerr

import (
	"errors"

	"bursavich.dev/errcode"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// An Error is an error with an explicit gRPC status.
type Error interface {
	GRPCStatus() *status.Status
	error
}

var errorCoder errcode.ErrorCoder = errcode.FromFunc(ErrorCode)

// ErrorCoder return the gRPC ErrorCoder.
func ErrorCoder() errcode.ErrorCoder {
	return errorCoder
}

// ErrorCode returns the gRPC code associated with the given error
// if it implements the gRPC Error interface.
func ErrorCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	gs, ok := err.(Error)
	if !ok && !errors.As(err, &gs) {
		return codes.Unknown
	}
	s := gs.GRPCStatus()
	if s == nil {
		// Error has status nil, which maps to codes.OK.
		// There is no sensible behavior for this, so we
		// turn it into an error with codes.Unknown.
		return codes.Unknown
	}
	return s.Code()
}
