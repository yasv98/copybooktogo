// Code generated by mockery v2.49.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// TestingT is an autogenerated mock type for the TestingT type
type TestingT struct {
	mock.Mock
}

type TestingT_Expecter struct {
	mock *mock.Mock
}

func (_m *TestingT) EXPECT() *TestingT_Expecter {
	return &TestingT_Expecter{mock: &_m.Mock}
}

// Errorf provides a mock function with given fields: format, args
func (_m *TestingT) Errorf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// TestingT_Errorf_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Errorf'
type TestingT_Errorf_Call struct {
	*mock.Call
}

// Errorf is a helper method to define mock.On call
//   - format string
//   - args ...interface{}
func (_e *TestingT_Expecter) Errorf(format interface{}, args ...interface{}) *TestingT_Errorf_Call {
	return &TestingT_Errorf_Call{Call: _e.mock.On("Errorf",
		append([]interface{}{format}, args...)...)}
}

func (_c *TestingT_Errorf_Call) Run(run func(format string, args ...interface{})) *TestingT_Errorf_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *TestingT_Errorf_Call) Return() *TestingT_Errorf_Call {
	_c.Call.Return()
	return _c
}

func (_c *TestingT_Errorf_Call) RunAndReturn(run func(string, ...interface{})) *TestingT_Errorf_Call {
	_c.Call.Return(run)
	return _c
}

// FailNow provides a mock function with given fields:
func (_m *TestingT) FailNow() {
	_m.Called()
}

// TestingT_FailNow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FailNow'
type TestingT_FailNow_Call struct {
	*mock.Call
}

// FailNow is a helper method to define mock.On call
func (_e *TestingT_Expecter) FailNow() *TestingT_FailNow_Call {
	return &TestingT_FailNow_Call{Call: _e.mock.On("FailNow")}
}

func (_c *TestingT_FailNow_Call) Run(run func()) *TestingT_FailNow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TestingT_FailNow_Call) Return() *TestingT_FailNow_Call {
	_c.Call.Return()
	return _c
}

func (_c *TestingT_FailNow_Call) RunAndReturn(run func()) *TestingT_FailNow_Call {
	_c.Call.Return(run)
	return _c
}

// NewTestingT creates a new instance of TestingT. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTestingT(t interface {
	mock.TestingT
	Cleanup(func())
}) *TestingT {
	mock := &TestingT{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
