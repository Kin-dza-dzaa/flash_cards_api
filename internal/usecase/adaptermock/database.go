// Code generated by mockery v2.20.0. DO NOT EDIT.

package adaptermock

import (
	context "context"

	entity "github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	mock "github.com/stretchr/testify/mock"
)

// Database is an autogenerated mock type for the Database type
type Database struct {
	mock.Mock
}

// AddTranslation provides a mock function with given fields: ctx, wordTrans
func (_m *Database) AddTranslation(ctx context.Context, wordTrans entity.WordTrans) error {
	ret := _m.Called(ctx, wordTrans)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.WordTrans) error); ok {
		r0 = rf(ctx, wordTrans)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddWordToCollection provides a mock function with given fields: ctx, collection
func (_m *Database) AddWordToCollection(ctx context.Context, collection entity.Collection) error {
	ret := _m.Called(ctx, collection)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Collection) error); ok {
		r0 = rf(ctx, collection)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteWordFromCollection provides a mock function with given fields: ctx, collection
func (_m *Database) DeleteWordFromCollection(ctx context.Context, collection entity.Collection) error {
	ret := _m.Called(ctx, collection)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Collection) error); ok {
		r0 = rf(ctx, collection)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserWords provides a mock function with given fields: ctx, collection
func (_m *Database) GetUserWords(ctx context.Context, collection entity.Collection) (*entity.UserWords, error) {
	ret := _m.Called(ctx, collection)

	var r0 *entity.UserWords
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Collection) (*entity.UserWords, error)); ok {
		return rf(ctx, collection)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.Collection) *entity.UserWords); ok {
		r0 = rf(ctx, collection)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.UserWords)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.Collection) error); ok {
		r1 = rf(ctx, collection)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsTransInDB provides a mock function with given fields: ctx, collection
func (_m *Database) IsTransInDB(ctx context.Context, collection entity.Collection) (bool, error) {
	ret := _m.Called(ctx, collection)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Collection) (bool, error)); ok {
		return rf(ctx, collection)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.Collection) bool); ok {
		r0 = rf(ctx, collection)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.Collection) error); ok {
		r1 = rf(ctx, collection)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsWordInCollection provides a mock function with given fields: ctx, collection
func (_m *Database) IsWordInCollection(ctx context.Context, collection entity.Collection) (bool, error) {
	ret := _m.Called(ctx, collection)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Collection) (bool, error)); ok {
		return rf(ctx, collection)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.Collection) bool); ok {
		r0 = rf(ctx, collection)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.Collection) error); ok {
		r1 = rf(ctx, collection)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateLearnInterval provides a mock function with given fields: ctx, collection
func (_m *Database) UpdateLearnInterval(ctx context.Context, collection entity.Collection) error {
	ret := _m.Called(ctx, collection)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Collection) error); ok {
		r0 = rf(ctx, collection)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewDatabase interface {
	mock.TestingT
	Cleanup(func())
}

// NewDatabase creates a new instance of database. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDatabase(t mockConstructorTestingTnewDatabase) *Database {
	mock := &Database{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
