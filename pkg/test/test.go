// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.

package test

//go:generate mockgen -destination=../../mock/test.go -package=mock github.com/hbocodelabs/infratest/pkg/test T
type T interface {
	Errorf(string, ...interface{})
	FailNow()
	Fail()
	Log(...interface{})
	Logf(string, ...interface{})
}
