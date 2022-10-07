// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	schedule "github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	mock "github.com/stretchr/testify/mock"
)

// IOpsgenieSchedule is an autogenerated mock type for the iOpsgenieSchedule type
type IOpsgenieSchedule struct {
	mock.Mock
}

type IOpsgenieSchedule_Expecter struct {
	mock *mock.Mock
}

func (_m *IOpsgenieSchedule) EXPECT() *IOpsgenieSchedule_Expecter {
	return &IOpsgenieSchedule_Expecter{mock: &_m.Mock}
}

// GetOnCalls provides a mock function with given fields: _a0, request
func (_m *IOpsgenieSchedule) GetOnCalls(_a0 context.Context, request *schedule.GetOnCallsRequest) (*schedule.GetOnCallsResult, error) {
	ret := _m.Called(_a0, request)

	var r0 *schedule.GetOnCallsResult
	if rf, ok := ret.Get(0).(func(context.Context, *schedule.GetOnCallsRequest) *schedule.GetOnCallsResult); ok {
		r0 = rf(_a0, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*schedule.GetOnCallsResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *schedule.GetOnCallsRequest) error); ok {
		r1 = rf(_a0, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IOpsgenieSchedule_GetOnCalls_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOnCalls'
type IOpsgenieSchedule_GetOnCalls_Call struct {
	*mock.Call
}

// GetOnCalls is a helper method to define mock.On call
//   - _a0 context.Context
//   - request *schedule.GetOnCallsRequest
func (_e *IOpsgenieSchedule_Expecter) GetOnCalls(_a0 interface{}, request interface{}) *IOpsgenieSchedule_GetOnCalls_Call {
	return &IOpsgenieSchedule_GetOnCalls_Call{Call: _e.mock.On("GetOnCalls", _a0, request)}
}

func (_c *IOpsgenieSchedule_GetOnCalls_Call) Run(run func(_a0 context.Context, request *schedule.GetOnCallsRequest)) *IOpsgenieSchedule_GetOnCalls_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*schedule.GetOnCallsRequest))
	})
	return _c
}

func (_c *IOpsgenieSchedule_GetOnCalls_Call) Return(_a0 *schedule.GetOnCallsResult, _a1 error) *IOpsgenieSchedule_GetOnCalls_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewIOpsgenieSchedule interface {
	mock.TestingT
	Cleanup(func())
}

// NewIOpsgenieSchedule creates a new instance of IOpsgenieSchedule. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIOpsgenieSchedule(t mockConstructorTestingTNewIOpsgenieSchedule) *IOpsgenieSchedule {
	mock := &IOpsgenieSchedule{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
