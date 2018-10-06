// Copyright 2013 Canonical Ltd.
// Copyright 2014 Cloudbase Solutions SRL
// Licensed under the LGPLv3, see LICENCE file for details.

// +build windows

package utils_test

import (
	gc "gopkg.in/check.v1"

	"github.com/juju/utils"
)

type windowsFileSuite struct {
}

var _ = gc.Suite(&windowsFileSuite{})

func (s *windowsFileSuite) TestMakeFileURL(c *gc.C) {
	var makeFileURLTests = []struct {
		in       string
		expected string
	}{{
		in:       "file://C:\\foo\\baz",
		expected: "file://C:/foo/baz",
	}, {
		in:       "C:\\foo\\baz",
		expected: "file://C:/foo/baz",
	}, {
		in:       "http://foo/baz",
		expected: "http://foo/baz",
	}, {
		in:       "file://C:/foo/baz",
		expected: "file://C:/foo/baz",
	}}

	for i, t := range makeFileURLTests {
		c.Logf("Test %d", i)
		c.Assert(utils.MakeFileURL(t.in), gc.Equals, t.expected)
	}
}
