package service

import (
	"context"
	"testing"

	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/service/mocks"
)

func TestStore(t *testing.T) {
	pwConn := mocks.NewPwRepositoryInterface(t)
	pwServiceConn := NewPassword(pwConn)

	testTable := []struct {
		name      string
		input     model.User
		context   context.Context
		someError error
	}{
		{
			name: "std input with ID=110",
			input: model.User{
				ID:       110,
				Password: "Abcd",
			},
			context:   context.TODO(),
			someError: nil,
		},
		{
			name: "std input with ID=113",
			input: model.User{
				ID:       113,
				Password: "S0pranO",
			},
			context:   context.TODO(),
			someError: nil,
		},
		{
			name: "std input with someError",
			input: model.User{
				ID:       133,
				Password: "adsd",
			},
			context:   context.TODO(),
			someError: errSome,
		},
	}
	for _, test := range testTable {
		pwCall := pwConn.EXPECT().Store(test.context, test.input).Return(test.someError)
		pwServiceConn.Store(test.context, test.input)
		pwConn.AssertExpectations(t)
		pwCall.Unset()
	}

}

func TestCompare(t *testing.T) {
	pwConn := mocks.NewPwRepositoryInterface(t)
	pwServiceConn := NewPassword(pwConn)

	testTable := []struct {
		name      string
		input     model.User
		output    bool
		context   context.Context
		someError error
	}{
		{
			name: "std input with ID=110 and correct password",
			input: model.User{
				ID:       110,
				Password: "Abcd",
			},
			output:    true,
			context:   context.TODO(),
			someError: nil,
		},
		{
			name: "std input with ID=113 and incorect password",
			input: model.User{
				ID:       113,
				Password: "S0pranO",
			},
			output:    false,
			context:   context.TODO(),
			someError: nil,
		},
		{
			name: "std input with someError",
			input: model.User{
				ID:       133,
				Password: "adsd",
			},
			output:    false,
			context:   context.TODO(),
			someError: errSome,
		},
	}
	for _, test := range testTable {
		pwCall := pwConn.EXPECT().Compare(test.context, test.input).Return(test.output, test.someError)
		pwServiceConn.Compare(test.context, test.input)
		pwConn.AssertExpectations(t)
		pwCall.Unset()
	}
}

func TestDeletePassword(t *testing.T) {
	pwConn := mocks.NewPwRepositoryInterface(t)
	pwServiceConn := NewPassword(pwConn)

	testTable := []struct {
		name      string
		input     int
		context   context.Context
		someError error
	}{
		{
			name:      "std input with ID=110",
			input:     110,
			context:   context.TODO(),
			someError: nil,
		},
		{
			name:      "std input with ID=113",
			input:     113,
			context:   context.TODO(),
			someError: nil,
		},
		{
			name:      "std input with someError",
			input:     133,
			context:   context.TODO(),
			someError: errSome,
		},
	}
	for _, test := range testTable {
		pwCall := pwConn.EXPECT().DeletePassword(test.context, test.input).Return(test.someError)
		pwServiceConn.DeletePassword(test.context, test.input)
		pwConn.AssertExpectations(t)
		pwCall.Unset()
	}

}
