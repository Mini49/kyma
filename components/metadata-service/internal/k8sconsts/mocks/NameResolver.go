// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// NameResolver is an autogenerated mock type for the NameResolver type
type NameResolver struct {
	mock.Mock
}

// ExtractServiceId provides a mock function with given fields: remoteEnvironment, host
func (_m *NameResolver) ExtractServiceId(remoteEnvironment string, host string) string {
	ret := _m.Called(remoteEnvironment, host)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(remoteEnvironment, host)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetGatewayUrl provides a mock function with given fields: remoteEnvironment, id
func (_m *NameResolver) GetGatewayUrl(remoteEnvironment string, id string) string {
	ret := _m.Called(remoteEnvironment, id)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(remoteEnvironment, id)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetResourceName provides a mock function with given fields: remoteEnvironment, id
func (_m *NameResolver) GetResourceName(remoteEnvironment string, id string) string {
	ret := _m.Called(remoteEnvironment, id)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(remoteEnvironment, id)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
