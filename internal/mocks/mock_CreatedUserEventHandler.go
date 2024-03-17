// Code generated by mockery v2.41.0. DO NOT EDIT.

package mocks

import (
	context "context"

	create "github.com/royhq/go-play-app/features/users/create"
	mock "github.com/stretchr/testify/mock"
)

// CreatedUserEventHandlerMock is an autogenerated mock type for the CreatedUserEventHandler type
type CreatedUserEventHandlerMock struct {
	mock.Mock
}

type CreatedUserEventHandlerMock_Expecter struct {
	mock *mock.Mock
}

func (_m *CreatedUserEventHandlerMock) EXPECT() *CreatedUserEventHandlerMock_Expecter {
	return &CreatedUserEventHandlerMock_Expecter{mock: &_m.Mock}
}

// Handle provides a mock function with given fields: _a0, _a1
func (_m *CreatedUserEventHandlerMock) Handle(_a0 context.Context, _a1 create.CreatedUserEvent) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Handle")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, create.CreatedUserEvent) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreatedUserEventHandlerMock_Handle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Handle'
type CreatedUserEventHandlerMock_Handle_Call struct {
	*mock.Call
}

// Handle is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 create.CreatedUserEvent
func (_e *CreatedUserEventHandlerMock_Expecter) Handle(_a0 interface{}, _a1 interface{}) *CreatedUserEventHandlerMock_Handle_Call {
	return &CreatedUserEventHandlerMock_Handle_Call{Call: _e.mock.On("Handle", _a0, _a1)}
}

func (_c *CreatedUserEventHandlerMock_Handle_Call) Run(run func(_a0 context.Context, _a1 create.CreatedUserEvent)) *CreatedUserEventHandlerMock_Handle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(create.CreatedUserEvent))
	})
	return _c
}

func (_c *CreatedUserEventHandlerMock_Handle_Call) Return(_a0 error) *CreatedUserEventHandlerMock_Handle_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CreatedUserEventHandlerMock_Handle_Call) RunAndReturn(run func(context.Context, create.CreatedUserEvent) error) *CreatedUserEventHandlerMock_Handle_Call {
	_c.Call.Return(run)
	return _c
}

// NewCreatedUserEventHandlerMock creates a new instance of CreatedUserEventHandlerMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCreatedUserEventHandlerMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *CreatedUserEventHandlerMock {
	mock := &CreatedUserEventHandlerMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
