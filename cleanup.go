// Copyright 2013, 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package testing

import (
  "os/exec"

  gc "launchpad.net/gocheck"
)

type CleanupFunc func(*gc.C)
type cleanupStack []CleanupFunc

// CleanupSuite adds the ability to add cleanup functions that are called
// during either test tear down or suite tear down depending on the method
// called.
type CleanupSuite struct {
  testStack  cleanupStack
  suiteStack cleanupStack
  setupSuite bool
}

func (s *CleanupSuite) SetUpSuite(c *gc.C) {
  s.suiteStack = nil
  s.setupSuite = true
}

func (s *CleanupSuite) TearDownSuite(c *gc.C) {
  s.callStack(c, s.suiteStack)
  s.setupSuite = false
}

func (s *CleanupSuite) SetUpTest(c *gc.C) {
  s.setupSuite = false
  s.testStack = nil
}

func (s *CleanupSuite) TearDownTest(c *gc.C) {
  s.callStack(c, s.testStack)
}

func (s *CleanupSuite) callStack(c *gc.C, stack cleanupStack) {
  for i := len(stack) - 1; i >= 0; i-- {
    stack[i](c)
  }
}

// AddCleanup pushes the cleanup function onto the stack of functions to be
// called during TearDownTest.
func (s *CleanupSuite) AddCleanup(cleanup CleanupFunc) {
  s.testStack = append(s.testStack, cleanup)
}

// AddSuiteCleanup pushes the cleanup function onto the stack of functions to
// be called during TearDownSuite.
func (s *CleanupSuite) AddSuiteCleanup(cleanup CleanupFunc) {
  s.suiteStack = append(s.suiteStack, cleanup)
}

// PatchEnvironment sets the environment variable 'name' the the value passed
// in. The old value is saved and returned to the original value at test tear
// down time using a cleanup function.
func (s *CleanupSuite) PatchEnvironment(name, value string) {
  restore := PatchEnvironment(name, value)
  if s.setupSuite {
    s.AddSuiteCleanup(func(*gc.C) { restore() })
  } else {
    s.AddCleanup(func(*gc.C) { restore() })
  }
}

// PatchEnvPathPrepend prepends the given path to the environment $PATH and restores the
// original path on test teardown.
func (s *CleanupSuite) PatchEnvPathPrepend(dir string) {
  restore := PatchEnvPathPrepend(dir)
  if s.setupSuite {
    s.AddSuiteCleanup(func(*gc.C) { restore() })
  } else {
    s.AddCleanup(func(*gc.C) { restore() })
  }
}

// PatchValue sets the 'dest' variable the the value passed in. The old value
// is saved and returned to the original value at test tear down time using a
// cleanup function. The value must be assignable to the element type of the
// destination.
func (s *CleanupSuite) PatchValue(dest, value interface{}) {
  restore := PatchValue(dest, value)
  if s.setupSuite {
    s.AddSuiteCleanup(func(*gc.C) { restore() })
  } else {
    s.AddCleanup(func(*gc.C) { restore() })
  }
}

// HookCommandOutput calls the package function of the same name to mock out
// the result of a particular comand execution, and will call the restore
// function on test teardown.
func (s *CleanupSuite) HookCommandOutput(
  outputFunc *func(cmd *exec.Cmd) ([]byte, error),
  output []byte,
  err error,
) <-chan *exec.Cmd {
  result, restore := HookCommandOutput(outputFunc, output, err)
  s.AddCleanup(func(*gc.C) { restore() })
  return result
}
