// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package teststore

import (
	"testing"

	"czarcoin.org/czarcoin/storage/testsuite"
)

func TestSuite(t *testing.T)      { testsuite.RunTests(t, New()) }
func BenchmarkSuite(b *testing.B) { testsuite.RunBenchmarks(b, New()) }
