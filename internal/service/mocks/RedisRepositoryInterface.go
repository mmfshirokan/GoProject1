// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/mmfshirokan/GoProject1/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// RedisRepositoryInterface is an autogenerated mock type for the RedisRepositoryInterface type
type RedisRepositoryInterface[object interface {
	*model.User | []*model.RefreshToken
}] struct {
	mock.Mock
}

type RedisRepositoryInterface_Expecter[object interface {
	*model.User | []*model.RefreshToken
}] struct {
	mock *mock.Mock
}

func (_m *RedisRepositoryInterface[object]) EXPECT() *RedisRepositoryInterface_Expecter[object] {
	return &RedisRepositoryInterface_Expecter[object]{mock: &_m.Mock}
}

// Add provides a mock function with given fields: ctx, key, obj
func (_m *RedisRepositoryInterface[object]) Add(ctx context.Context, key string, obj object) error {
	ret := _m.Called(ctx, key, obj)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, object) error); ok {
		r0 = rf(ctx, key, obj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RedisRepositoryInterface_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type RedisRepositoryInterface_Add_Call[object interface {
	*model.User | []*model.RefreshToken
}] struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - obj object
func (_e *RedisRepositoryInterface_Expecter[object]) Add(ctx interface{}, key interface{}, obj interface{}) *RedisRepositoryInterface_Add_Call[object] {
	return &RedisRepositoryInterface_Add_Call[object]{Call: _e.mock.On("Add", ctx, key, obj)}
}

func (_c *RedisRepositoryInterface_Add_Call[object]) Run(run func(ctx context.Context, key string, obj object)) *RedisRepositoryInterface_Add_Call[object] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(object))
	})
	return _c
}

func (_c *RedisRepositoryInterface_Add_Call[object]) Return(_a0 error) *RedisRepositoryInterface_Add_Call[object] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RedisRepositoryInterface_Add_Call[object]) RunAndReturn(run func(context.Context, string, object) error) *RedisRepositoryInterface_Add_Call[object] {
	_c.Call.Return(run)
	return _c
}

// NewRedisRepositoryInterface creates a new instance of RedisRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRedisRepositoryInterface[object interface {
	*model.User | []*model.RefreshToken
}](t interface {
	mock.TestingT
	Cleanup(func())
}) *RedisRepositoryInterface[object] {
	mock := &RedisRepositoryInterface[object]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}