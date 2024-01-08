// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/mmfshirokan/GoProject1/internal/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// AuthRepositoryInterface is an autogenerated mock type for the AuthRepositoryInterface type
type AuthRepositoryInterface struct {
	mock.Mock
}

type AuthRepositoryInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *AuthRepositoryInterface) EXPECT() *AuthRepositoryInterface_Expecter {
	return &AuthRepositoryInterface_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, token
func (_m *AuthRepositoryInterface) Create(ctx context.Context, token *model.RefreshToken) error {
	ret := _m.Called(ctx, token)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.RefreshToken) error); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthRepositoryInterface_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type AuthRepositoryInterface_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - token *model.RefreshToken
func (_e *AuthRepositoryInterface_Expecter) Create(ctx interface{}, token interface{}) *AuthRepositoryInterface_Create_Call {
	return &AuthRepositoryInterface_Create_Call{Call: _e.mock.On("Create", ctx, token)}
}

func (_c *AuthRepositoryInterface_Create_Call) Run(run func(ctx context.Context, token *model.RefreshToken)) *AuthRepositoryInterface_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.RefreshToken))
	})
	return _c
}

func (_c *AuthRepositoryInterface_Create_Call) Return(_a0 error) *AuthRepositoryInterface_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AuthRepositoryInterface_Create_Call) RunAndReturn(run func(context.Context, *model.RefreshToken) error) *AuthRepositoryInterface_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, id
func (_m *AuthRepositoryInterface) Delete(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthRepositoryInterface_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type AuthRepositoryInterface_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
func (_e *AuthRepositoryInterface_Expecter) Delete(ctx interface{}, id interface{}) *AuthRepositoryInterface_Delete_Call {
	return &AuthRepositoryInterface_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *AuthRepositoryInterface_Delete_Call) Run(run func(ctx context.Context, id uuid.UUID)) *AuthRepositoryInterface_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *AuthRepositoryInterface_Delete_Call) Return(_a0 error) *AuthRepositoryInterface_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AuthRepositoryInterface_Delete_Call) RunAndReturn(run func(context.Context, uuid.UUID) error) *AuthRepositoryInterface_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetByUserID provides a mock function with given fields: ctx, userID
func (_m *AuthRepositoryInterface) GetByUserID(ctx context.Context, userID int) ([]*model.RefreshToken, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetByUserID")
	}

	var r0 []*model.RefreshToken
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) ([]*model.RefreshToken, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) []*model.RefreshToken); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.RefreshToken)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthRepositoryInterface_GetByUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByUserID'
type AuthRepositoryInterface_GetByUserID_Call struct {
	*mock.Call
}

// GetByUserID is a helper method to define mock.On call
//   - ctx context.Context
//   - userID int
func (_e *AuthRepositoryInterface_Expecter) GetByUserID(ctx interface{}, userID interface{}) *AuthRepositoryInterface_GetByUserID_Call {
	return &AuthRepositoryInterface_GetByUserID_Call{Call: _e.mock.On("GetByUserID", ctx, userID)}
}

func (_c *AuthRepositoryInterface_GetByUserID_Call) Run(run func(ctx context.Context, userID int)) *AuthRepositoryInterface_GetByUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *AuthRepositoryInterface_GetByUserID_Call) Return(_a0 []*model.RefreshToken, _a1 error) *AuthRepositoryInterface_GetByUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthRepositoryInterface_GetByUserID_Call) RunAndReturn(run func(context.Context, int) ([]*model.RefreshToken, error)) *AuthRepositoryInterface_GetByUserID_Call {
	_c.Call.Return(run)
	return _c
}

// NewAuthRepositoryInterface creates a new instance of AuthRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthRepositoryInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuthRepositoryInterface {
	mock := &AuthRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
