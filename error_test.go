// SPDX-License-Identifier: MIT
//
// Copyright 2025 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

package errcode

import (
	"reflect"
	"slices"
	"testing"
)

func TestCompact(t *testing.T) {
	want := ErrorCoders{
		CodedErrorCoder(),
		ContextErrorCoder(),
		FileSystemErrorCoder(),
	}
	got := Compact(slices.Repeat(want, 3))
	if !reflect.DeepEqual(got, want) {
		t.Fail()
	}
}
