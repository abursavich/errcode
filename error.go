// SPDX-License-Identifier: MIT
//
// Copyright 2025 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

// Package errcode provides the ability to associate a gRPC status code with am error.
package errcode

import (
	"context"
	"errors"
	"io/fs"
	"slices"

	"google.golang.org/grpc/codes"
)

// An Error is an error with an explicit code.
type Error interface {
	Code() codes.Code
	error
}

// New wraps the given error and adds an explicit code.
func New(code codes.Code, err error) error {
	return &codedError{code, err}
}

type codedError struct {
	code codes.Code
	err  error
}

func (ce *codedError) Code() codes.Code { return ce.code }
func (ce *codedError) Error() string    { return ce.err.Error() }
func (ce *codedError) Unwrap() error    { return ce.err }

// An ErrorCoder returns the code of an error.
type ErrorCoder interface {
	// ErrorCode returns the code of an error.
	// If the error is nil, it must return OK.
	// If the code cannot be determined, it must return Unknown.
	ErrorCode(error) codes.Code
}

type errorCoderFn struct {
	fn func(error) codes.Code
}

// FromFunc returns an ErrorCoder from a function.
func FromFunc(fn func(error) codes.Code) ErrorCoder {
	return &errorCoderFn{fn}
}

func (e *errorCoderFn) ErrorCode(err error) codes.Code {
	return e.fn(err)
}

// ErrorCoders is an ErrorCoder that combines other ErrorCoders.
type ErrorCoders []ErrorCoder

func (s ErrorCoders) ErrorCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	for _, v := range s {
		if c := v.ErrorCode(err); c != codes.Unknown {
			return c
		}
	}
	return codes.Unknown
}

// Compact flattens and dedupes ErrorCoders.
func Compact(coders ...ErrorCoder) ErrorCoders {
	return compact(nil, coders...)
}

func compact(slice ErrorCoders, elems ...ErrorCoder) ErrorCoders {
	for _, elem := range elems {
		if list, ok := elem.(ErrorCoders); ok {
			slice = compact(slice, list...)
			continue
		}
		if !contains(slice, elem) {
			slice = append(slice, elem)
		}
	}
	return slice
}

func contains(slice ErrorCoders, elem ErrorCoder) bool {
	defer func() { _ = recover() }()
	return slices.Contains(slice, elem)
}

var codedErrorCoder ErrorCoder = FromFunc(codedErrorCode)

// CodedErrorCoder returns an ErrorCoder that handles CodedErrors.
func CodedErrorCoder() ErrorCoder {
	return codedErrorCoder
}

func codedErrorCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	var e Error
	if errors.As(err, &e) {
		return e.Code()
	}
	return codes.Unknown
}

var contextErrorCoder ErrorCoder = FromFunc(contextErrorCode)

// ContextErrorCoder returns an ErrorCoder that handles context errors.
func ContextErrorCoder() ErrorCoder {
	return contextErrorCoder
}

func contextErrorCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return codes.DeadlineExceeded
	}
	if errors.Is(err, context.Canceled) {
		return codes.Canceled
	}
	return codes.Unknown
}

var fsErrorCoder ErrorCoder = FromFunc(fsErrorCode)

// FileSystemErrorCoder returns an ErrorCoder that handles fs errors.
func FileSystemErrorCoder() ErrorCoder {
	return fsErrorCoder
}

func fsErrorCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if errors.Is(err, fs.ErrExist) {
		return codes.AlreadyExists
	}
	if errors.Is(err, fs.ErrNotExist) {
		return codes.NotFound
	}
	if errors.Is(err, fs.ErrPermission) {
		return codes.PermissionDenied
	}
	if errors.Is(err, fs.ErrInvalid) {
		return codes.InvalidArgument
	}
	return codes.Unknown
}
