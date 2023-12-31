// Code generated by mockery v2.33.2. DO NOT EDIT.

package keyservice

import mock "github.com/stretchr/testify/mock"

// MockKeyService is an autogenerated mock type for the KeyService type
type MockKeyService struct {
	mock.Mock
}

type MockKeyService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockKeyService) EXPECT() *MockKeyService_Expecter {
	return &MockKeyService_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: _a0
func (_m *MockKeyService) Delete(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockKeyService_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockKeyService_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - _a0 string
func (_e *MockKeyService_Expecter) Delete(_a0 interface{}) *MockKeyService_Delete_Call {
	return &MockKeyService_Delete_Call{Call: _e.mock.On("Delete", _a0)}
}

func (_c *MockKeyService_Delete_Call) Run(run func(_a0 string)) *MockKeyService_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockKeyService_Delete_Call) Return(_a0 error) *MockKeyService_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockKeyService_Delete_Call) RunAndReturn(run func(string) error) *MockKeyService_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: _a0
func (_m *MockKeyService) Get(_a0 string) (string, error) {
	ret := _m.Called(_a0)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockKeyService_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockKeyService_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - _a0 string
func (_e *MockKeyService_Expecter) Get(_a0 interface{}) *MockKeyService_Get_Call {
	return &MockKeyService_Get_Call{Call: _e.mock.On("Get", _a0)}
}

func (_c *MockKeyService_Get_Call) Run(run func(_a0 string)) *MockKeyService_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockKeyService_Get_Call) Return(_a0 string, _a1 error) *MockKeyService_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockKeyService_Get_Call) RunAndReturn(run func(string) (string, error)) *MockKeyService_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Put provides a mock function with given fields: _a0, _a1
func (_m *MockKeyService) Put(_a0 string, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockKeyService_Put_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Put'
type MockKeyService_Put_Call struct {
	*mock.Call
}

// Put is a helper method to define mock.On call
//   - _a0 string
//   - _a1 string
func (_e *MockKeyService_Expecter) Put(_a0 interface{}, _a1 interface{}) *MockKeyService_Put_Call {
	return &MockKeyService_Put_Call{Call: _e.mock.On("Put", _a0, _a1)}
}

func (_c *MockKeyService_Put_Call) Run(run func(_a0 string, _a1 string)) *MockKeyService_Put_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockKeyService_Put_Call) Return(_a0 error) *MockKeyService_Put_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockKeyService_Put_Call) RunAndReturn(run func(string, string) error) *MockKeyService_Put_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockKeyService creates a new instance of MockKeyService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockKeyService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockKeyService {
	mock := &MockKeyService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
