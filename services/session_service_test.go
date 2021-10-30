package services

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"testing"
)

func TestFetchSession(t *testing.T) {
	tests := []struct {
		SessionId     string
		Prepare       func(client *redismock.ClientMock)
		Expected      string
		ExpectedError error
	}{
		{
			SessionId: "ididntchosethisshinylife",
			Prepare: func(client *redismock.ClientMock) {
				(*client).ExpectGet("ididntchosethisshinylife").SetVal("aliparlakci")
			},
			Expected:      "aliparlakci",
			ExpectedError: nil,
		}, {
			SessionId: "ididntchosethisshinylife",
			Prepare: func(client *redismock.ClientMock) {
				(*client).ExpectGet("ididntchosethisshinylife").SetErr(redis.Nil)
			},
			Expected:      "",
			ExpectedError: ErrNoSession,
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("[%v]", i)
		t.Run(testName, func(t *testing.T) {
			db, mock := redismock.NewClientMock()

			tt.Prepare(&mock)

			service := SessionService{Store: db}
			result, err := service.FetchSession(context.Background(), tt.SessionId)

			if err != tt.ExpectedError {
				t.Errorf("want %v, got %v", tt.ExpectedError, err)
			}
			if result != tt.Expected {
				t.Errorf("want %v, got %v", tt.Expected, result)
			}
		})
	}
}

func TestCreateSession(t *testing.T) {
	tests := []struct {
		Data    string
		Prepare func(client *redismock.ClientMock)
		ExpectedError error
	}{
		{
			Data: "aliparlakci",
			Prepare: func(client *redismock.ClientMock) {
				(*client).Regexp().ExpectSet(`(.)*`, "aliparlakci", 0).SetVal("")
			},
			ExpectedError: nil,
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("[%v]", i)
		t.Run(testName, func(t *testing.T) {
			db, mock := redismock.NewClientMock()

			tt.Prepare(&mock)

			service := SessionService{Store: db}
			_, err := service.CreateSession(context.Background(), tt.Data)

			if err != tt.ExpectedError {
				t.Errorf("want %v, got %v", tt.ExpectedError, err)
			}
		})
	}
}

func TestRevokeSession(t *testing.T) {
	tests := []struct {
		Data    string
		Prepare func(client *redismock.ClientMock)
		ExpectedError error
	}{
		{
			Data: "someuuid",
			Prepare: func(client *redismock.ClientMock) {
				(*client).Regexp().ExpectDel("someuuid").SetErr(redis.Nil)
			},
			ExpectedError: redis.Nil,
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("[%v]", i)
		t.Run(testName, func(t *testing.T) {
			db, mock := redismock.NewClientMock()

			tt.Prepare(&mock)

			service := SessionService{Store: db}
			err := service.RevokeSession(context.Background(), tt.Data)

			if err != tt.ExpectedError {
				t.Errorf("want %v, got %v", tt.ExpectedError, err)
			}
		})
	}
}
