// SPDX-License-Identifier: MIT
//
// Copyright 2025 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

// Package googleapierr provides the ability to extract the status code from Google API errors.
package googleapierr

import (
	"errors"

	"bursavich.dev/errcode"
	"bursavich.dev/errcode/grpcerr"
	"bursavich.dev/errcode/httperr"
	"google.golang.org/api/googleapi"
	"google.golang.org/grpc/codes"
)

var errorCoder errcode.ErrorCoder = errcode.ErrorCoders{
	// NOTE: github.com/googleapis/gax-go/v2/apierror.APIError implements gRPC and HTTP error interfaces.
	// These are prefered over *google.golang.org/api/googleapi.Error which only include an HTTP code.
	grpcerr.ErrorCoder(),
	httperr.ErrorCoder(),
	errcode.FromFunc(googleAPIErrorCode),
}

func ErrorCoder() errcode.ErrorCoder {
	return errorCoder
}

func ErrorCode(err error) codes.Code {
	return errorCoder.ErrorCode(err)
}

func googleAPIErrorCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if ge, ok := err.(*googleapi.Error); ok || errors.As(err, &ge) {
		return httperr.ToGRPC(ge.Code)
	}
	return codes.Unknown
}
