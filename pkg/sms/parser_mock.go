// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package sms

import mock "github.com/stretchr/testify/mock"

// MockParser is an autogenerated mock type for the Parser type
type MockParser struct {
	mock.Mock
}

// Parse provides a mock function with given fields: text
func (_m *MockParser) Parse(text string) (Message, error) {
	ret := _m.Called(text)

	var r0 Message
	if rf, ok := ret.Get(0).(func(string) Message); ok {
		r0 = rf(text)
	} else {
		r0 = ret.Get(0).(Message)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(text)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}