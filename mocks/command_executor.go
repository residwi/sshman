// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package mocks

import (
	mock "github.com/stretchr/testify/mock"
)

// NewMockCommandExecutor creates a new instance of MockCommandExecutor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCommandExecutor(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCommandExecutor {
	mock := &MockCommandExecutor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockCommandExecutor is an autogenerated mock type for the CommandExecutor type
type MockCommandExecutor struct {
	mock.Mock
}

type MockCommandExecutor_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCommandExecutor) EXPECT() *MockCommandExecutor_Expecter {
	return &MockCommandExecutor_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function for the type MockCommandExecutor
func (_mock *MockCommandExecutor) Execute(name string, args ...string) error {
	var tmpRet mock.Arguments
	if len(args) > 0 {
		tmpRet = _mock.Called(name, args)
	} else {
		tmpRet = _mock.Called(name)
	}
	ret := tmpRet

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(string, ...string) error); ok {
		r0 = returnFunc(name, args...)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockCommandExecutor_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockCommandExecutor_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - name string
//   - args ...string
func (_e *MockCommandExecutor_Expecter) Execute(name interface{}, args ...interface{}) *MockCommandExecutor_Execute_Call {
	return &MockCommandExecutor_Execute_Call{Call: _e.mock.On("Execute",
		append([]interface{}{name}, args...)...)}
}

func (_c *MockCommandExecutor_Execute_Call) Run(run func(name string, args ...string)) *MockCommandExecutor_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 string
		if args[0] != nil {
			arg0 = args[0].(string)
		}
		var arg1 []string
		var variadicArgs []string
		if len(args) > 1 {
			variadicArgs = args[1].([]string)
		}
		arg1 = variadicArgs
		run(
			arg0,
			arg1...,
		)
	})
	return _c
}

func (_c *MockCommandExecutor_Execute_Call) Return(err error) *MockCommandExecutor_Execute_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockCommandExecutor_Execute_Call) RunAndReturn(run func(name string, args ...string) error) *MockCommandExecutor_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// ExecuteWithOutput provides a mock function for the type MockCommandExecutor
func (_mock *MockCommandExecutor) ExecuteWithOutput(name string, args ...string) ([]byte, error) {
	var tmpRet mock.Arguments
	if len(args) > 0 {
		tmpRet = _mock.Called(name, args)
	} else {
		tmpRet = _mock.Called(name)
	}
	ret := tmpRet

	if len(ret) == 0 {
		panic("no return value specified for ExecuteWithOutput")
	}

	var r0 []byte
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(string, ...string) ([]byte, error)); ok {
		return returnFunc(name, args...)
	}
	if returnFunc, ok := ret.Get(0).(func(string, ...string) []byte); ok {
		r0 = returnFunc(name, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(string, ...string) error); ok {
		r1 = returnFunc(name, args...)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockCommandExecutor_ExecuteWithOutput_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteWithOutput'
type MockCommandExecutor_ExecuteWithOutput_Call struct {
	*mock.Call
}

// ExecuteWithOutput is a helper method to define mock.On call
//   - name string
//   - args ...string
func (_e *MockCommandExecutor_Expecter) ExecuteWithOutput(name interface{}, args ...interface{}) *MockCommandExecutor_ExecuteWithOutput_Call {
	return &MockCommandExecutor_ExecuteWithOutput_Call{Call: _e.mock.On("ExecuteWithOutput",
		append([]interface{}{name}, args...)...)}
}

func (_c *MockCommandExecutor_ExecuteWithOutput_Call) Run(run func(name string, args ...string)) *MockCommandExecutor_ExecuteWithOutput_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 string
		if args[0] != nil {
			arg0 = args[0].(string)
		}
		var arg1 []string
		var variadicArgs []string
		if len(args) > 1 {
			variadicArgs = args[1].([]string)
		}
		arg1 = variadicArgs
		run(
			arg0,
			arg1...,
		)
	})
	return _c
}

func (_c *MockCommandExecutor_ExecuteWithOutput_Call) Return(bytes []byte, err error) *MockCommandExecutor_ExecuteWithOutput_Call {
	_c.Call.Return(bytes, err)
	return _c
}

func (_c *MockCommandExecutor_ExecuteWithOutput_Call) RunAndReturn(run func(name string, args ...string) ([]byte, error)) *MockCommandExecutor_ExecuteWithOutput_Call {
	_c.Call.Return(run)
	return _c
}

// ExecuteWithoutOutput provides a mock function for the type MockCommandExecutor
func (_mock *MockCommandExecutor) ExecuteWithoutOutput(name string, args ...string) error {
	var tmpRet mock.Arguments
	if len(args) > 0 {
		tmpRet = _mock.Called(name, args)
	} else {
		tmpRet = _mock.Called(name)
	}
	ret := tmpRet

	if len(ret) == 0 {
		panic("no return value specified for ExecuteWithoutOutput")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(string, ...string) error); ok {
		r0 = returnFunc(name, args...)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockCommandExecutor_ExecuteWithoutOutput_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteWithoutOutput'
type MockCommandExecutor_ExecuteWithoutOutput_Call struct {
	*mock.Call
}

// ExecuteWithoutOutput is a helper method to define mock.On call
//   - name string
//   - args ...string
func (_e *MockCommandExecutor_Expecter) ExecuteWithoutOutput(name interface{}, args ...interface{}) *MockCommandExecutor_ExecuteWithoutOutput_Call {
	return &MockCommandExecutor_ExecuteWithoutOutput_Call{Call: _e.mock.On("ExecuteWithoutOutput",
		append([]interface{}{name}, args...)...)}
}

func (_c *MockCommandExecutor_ExecuteWithoutOutput_Call) Run(run func(name string, args ...string)) *MockCommandExecutor_ExecuteWithoutOutput_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 string
		if args[0] != nil {
			arg0 = args[0].(string)
		}
		var arg1 []string
		var variadicArgs []string
		if len(args) > 1 {
			variadicArgs = args[1].([]string)
		}
		arg1 = variadicArgs
		run(
			arg0,
			arg1...,
		)
	})
	return _c
}

func (_c *MockCommandExecutor_ExecuteWithoutOutput_Call) Return(err error) *MockCommandExecutor_ExecuteWithoutOutput_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockCommandExecutor_ExecuteWithoutOutput_Call) RunAndReturn(run func(name string, args ...string) error) *MockCommandExecutor_ExecuteWithoutOutput_Call {
	_c.Call.Return(run)
	return _c
}
