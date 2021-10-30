package handlers

import (
	"errors"
	"fmt"
	"github.com/aliparlakci/armut-backend-assessment/common"
	"github.com/aliparlakci/armut-backend-assessment/mocks"
	"github.com/aliparlakci/armut-backend-assessment/models"
	"github.com/aliparlakci/armut-backend-assessment/services"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetAllMessages(t *testing.T) {
	tests := []struct {
		Prepare      func(getter *mocks.MockMessageGetter)
		ExpectedCode int
		ExpectedBody gin.H
	}{
		{
			Prepare: func(getter *mocks.MockMessageGetter) {
				getter.EXPECT().GetAllMessages(gomock.Any(), "johndoe").Return([]models.Message{
					{
						ID:     primitive.ObjectID{},
						From:   "iskralawrence",
						To:     "aliparlakci",
						Body:   "hey cutie :)",
						IsRead: false,
						SendAt: time.Time{},
					},
				}, nil).MinTimes(1)
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: gin.H{"result": []gin.H{
				{
					"id":      "000000000000000000000000",
					"from":    "iskralawrence",
					"to":      "aliparlakci",
					"body":    "hey cutie :)",
					"is_read": false,
					"send_at": time.Time{},
				},
			}},
		}, {
			Prepare: func(getter *mocks.MockMessageGetter) {
				getter.EXPECT().GetAllMessages(gomock.Any(), "johndoe").Return(nil, errors.New("")).MinTimes(1)
			},
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: gin.H{},
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("[%v]", i)
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockedMessageGetter := mocks.NewMockMessageGetter(ctrl)

			if tt.Prepare != nil {
				tt.Prepare(mockedMessageGetter)
			}

			recorder := httptest.NewRecorder()
			_, r := gin.CreateTestContext(recorder)
			r.Use(func(c *gin.Context) {
				c.Set("user", models.User{Username: "johndoe"})
			})
			r.GET("/api/messages", GetAllMessages(mockedMessageGetter))

			request, err := http.NewRequest(http.MethodGet, "/api/messages", nil)

			if err != nil {
				t.Fatal(err)
			}

			r.ServeHTTP(recorder, request)

			if bodyAssertion, err := common.IsBodyEqual(tt.ExpectedBody, recorder.Result().Body); err != nil {
				t.Fatal(err)
			} else if !bodyAssertion {
				t.Errorf("response bodies don't match")
			}

			if recorder.Result().StatusCode != tt.ExpectedCode {
				t.Errorf("want %v, got %v", tt.ExpectedCode, recorder.Result().StatusCode)
			}
		})
	}
}

func TestGetNewMessages(t *testing.T) {
	tests := []struct {
		Prepare      func(getter *mocks.MockMessageGetter)
		ExpectedCode int
		ExpectedBody gin.H
	}{
		{
			Prepare: func(getter *mocks.MockMessageGetter) {
				getter.EXPECT().GetNewMessages(gomock.Any(), "johndoe").Return([]models.Message{
					{
						ID:     primitive.ObjectID{},
						From:   "iskralawrence",
						To:     "aliparlakci",
						Body:   "hey cutie :)",
						IsRead: false,
						SendAt: time.Time{},
					},
				}, nil).MinTimes(1)
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: gin.H{"result": []gin.H{
				{
					"id":      "000000000000000000000000",
					"from":    "iskralawrence",
					"to":      "aliparlakci",
					"body":    "hey cutie :)",
					"is_read": false,
					"send_at": time.Time{},
				},
			}},
		}, {
			Prepare: func(getter *mocks.MockMessageGetter) {
				getter.EXPECT().GetNewMessages(gomock.Any(), "johndoe").Return(nil, errors.New("")).MinTimes(1)
			},
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: gin.H{},
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("[%v]", i)
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockedMessageGetter := mocks.NewMockMessageGetter(ctrl)

			if tt.Prepare != nil {
				tt.Prepare(mockedMessageGetter)
			}

			recorder := httptest.NewRecorder()
			_, r := gin.CreateTestContext(recorder)
			r.Use(func(c *gin.Context) {
				c.Set("user", models.User{Username: "johndoe"})
			})
			r.GET("/api/messages/new", GetNewMessages(mockedMessageGetter))

			request, err := http.NewRequest(http.MethodGet, "/api/messages/new", nil)

			if err != nil {
				t.Fatal(err)
			}

			r.ServeHTTP(recorder, request)

			if bodyAssertion, err := common.IsBodyEqual(tt.ExpectedBody, recorder.Result().Body); err != nil {
				t.Fatal(err)
			} else if !bodyAssertion {
				t.Errorf("response bodies don't match")
			}

			if recorder.Result().StatusCode != tt.ExpectedCode {
				t.Errorf("want %v, got %v", tt.ExpectedCode, recorder.Result().StatusCode)
			}
		})
	}
}

func TestCheckNewMessages(t *testing.T) {
	tests := []struct {
		Prepare      func(getter *mocks.MockMessageGetter)
		ExpectedCode int
		ExpectedBody gin.H
	}{
		{
			Prepare: func(getter *mocks.MockMessageGetter) {
				getter.EXPECT().CheckNewMessages(gomock.Any(), "johndoe").Return(5, nil).MinTimes(1)
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: gin.H{"result": 5},
		}, {
			Prepare: func(getter *mocks.MockMessageGetter) {
				getter.EXPECT().CheckNewMessages(gomock.Any(), "johndoe").Return(0, errors.New("")).MinTimes(1)
			},
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: gin.H{},
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("[%v]", i)
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockedMessageGetter := mocks.NewMockMessageGetter(ctrl)

			if tt.Prepare != nil {
				tt.Prepare(mockedMessageGetter)
			}

			recorder := httptest.NewRecorder()
			_, r := gin.CreateTestContext(recorder)
			r.Use(func(c *gin.Context) {
				c.Set("user", models.User{Username: "johndoe"})
			})
			r.GET("/api/messages/check", CheckNewMessages(mockedMessageGetter))

			request, err := http.NewRequest(http.MethodGet, "/api/messages/check", nil)

			if err != nil {
				t.Fatal(err)
			}

			r.ServeHTTP(recorder, request)

			if bodyAssertion, err := common.IsBodyEqual(tt.ExpectedBody, recorder.Result().Body); err != nil {
				t.Fatal(err)
			} else if !bodyAssertion {
				t.Errorf("response bodies don't match")
			}

			if recorder.Result().StatusCode != tt.ExpectedCode {
				t.Errorf("want %v, got %v", tt.ExpectedCode, recorder.Result().StatusCode)
			}
		})
	}
}

func TestSendMessage(t *testing.T) {
	tests := []struct {
		Body         multipart.Form
		Prepare      func(sender *mocks.MockMessageSender)
		ExpectedCode int
		ExpectedBody gin.H
	}{
		{
			Body: multipart.Form{
				Value: map[string][]string{
					"to":   {"ozkanugur"},
					"body": {"ozkan selam, mazhar ben. bugun biraz asabiydim ama mazaretim var. aksam anlatirim"},
				},
			},
			Prepare: func(sender *mocks.MockMessageSender) {
				sender.EXPECT().SendMessage(gomock.Any(), "ozkan selam, mazhar ben. bugun biraz asabiydim ama mazaretim var. aksam anlatirim", "mazhar", "ozkanugur").Return("", nil).MinTimes(1)
			},
			ExpectedCode: http.StatusCreated,
			ExpectedBody: gin.H{"result": "message is successfully sent"},
		}, {
			Body: multipart.Form{
				Value: map[string][]string{
					"body": {"keske arkadasim olsa da ben de ozkan abi gibi mesajlassam..."},
				},
			},
			Prepare: func(sender *mocks.MockMessageSender) {
				sender.EXPECT().SendMessage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: gin.H{},
		}, {
			Body: multipart.Form{
				Value: map[string][]string{
					"to": {"fuatguner"},
				},
			},
			Prepare: func(sender *mocks.MockMessageSender) {
				sender.EXPECT().SendMessage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: gin.H{},
		}, {
			Body: multipart.Form{
				Value: map[string][]string{
					"to":   {"tarkan"},
					"body": {"tarkanla mesajlasmak bu kadar kolay miymis yav"},
				},
			},
			Prepare: func(sender *mocks.MockMessageSender) {
				sender.EXPECT().SendMessage(gomock.Any(), "tarkanla mesajlasmak bu kadar kolay miymis yav", "mazhar", "tarkan").Return(
					"",
					services.ErrNoUser,
				).MinTimes(1)
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: gin.H{"result": "user does not exist"},
		}, {
			Body: multipart.Form{
				Value: map[string][]string{
					"to":   {"tarkan"},
					"body": {"tarkanla mesajlasmak bu kadar kolay miymis yav"},
				}},
			Prepare: func(sender *mocks.MockMessageSender) {
				sender.EXPECT().SendMessage(gomock.Any(), "tarkanla mesajlasmak bu kadar kolay miymis yav", "mazhar", "tarkan").Return("", errors.New("")).MinTimes(1)
			},
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: gin.H{},
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("[%v]", i)
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockedMessageSender := mocks.NewMockMessageSender(ctrl)

			if tt.Prepare != nil {
				tt.Prepare(mockedMessageSender)
			}

			recorder := httptest.NewRecorder()
			_, r := gin.CreateTestContext(recorder)
			r.Use(func(c *gin.Context) {
				c.Set("user", models.User{Username: "mazhar"})
			})
			r.POST("/api/messages/send", SendMessage(mockedMessageSender))

			request, err := http.NewRequest(http.MethodPost, "/api/messages/send", nil)
			request.MultipartForm = &tt.Body
			request.Header.Set("Content-Type", "multipart/form-data")

			if err != nil {
				t.Fatal(err)
			}

			r.ServeHTTP(recorder, request)

			if bodyAssertion, err := common.IsBodyEqual(tt.ExpectedBody, recorder.Result().Body); err != nil {
				t.Fatal(err)
			} else if !bodyAssertion {
				t.Errorf("response bodies don't match")
			}

			if recorder.Result().StatusCode != tt.ExpectedCode {
				t.Errorf("want %v, got %v", tt.ExpectedCode, recorder.Result().StatusCode)
			}
		})
	}
}
