// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import kube "github.com/olegsu/iris/pkg/kube"
import mock "github.com/stretchr/testify/mock"

// Kube is an autogenerated mock type for the Kube type
type Kube struct {
	mock.Mock
}

// ResourceByLabelsExist provides a mock function with given fields: _a0, _a1
func (_m *Kube) ResourceByLabelsExist(_a0 interface{}, _a1 map[string]string) (bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(interface{}, map[string]string) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}, map[string]string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetIRISConfigmap provides a mock function with given fields: _a0, _a1
func (_m *Kube) GetIRISConfigmap(_a0 string, _a1 string) ([]byte, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, string) []byte); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Watch provides a mock function with given fields: _a0
func (_m *Kube) Watch(_a0 kube.WatchFn) {
	_m.Called(_a0)
}
