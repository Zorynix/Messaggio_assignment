package v1_test

import (
	"context"
	"encoding/json"
	"messagio_testsuite/internal/entity"
	v1 "messagio_testsuite/internal/routes/http/v1"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) CreateMessage(ctx context.Context, content string) (uuid.UUID, error) {
	args := m.Called(ctx, content)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockMessageService) GetMessageById(ctx context.Context, id uuid.UUID) (entity.Message, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Message), args.Error(1)
}

func (m *MockMessageService) GetMessages(ctx context.Context) ([]entity.Message, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Message), args.Error(1)
}

func (m *MockMessageService) MarkMessageAsProcessed(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMessageService) GetProcessedMessagesStats(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Get(0).(int), args.Error(1)
}

func setup() (*echo.Echo, *MockMessageService, *v1.MessageRoutes) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	mockService := new(MockMessageService)
	routes := &v1.MessageRoutes{
		MessageService: mockService,
	}
	v1.NewMessageRoutes(e.Group("/"), mockService)
	return e, mockService, routes
}

func TestCreateMessage(t *testing.T) {
	e, mockService, routes := setup()

	reqBody := `{"message": "Hello, world!"}`
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("CreateMessage", mock.Anything, "Hello, world!").Return(uuid.New(), nil)

	if assert.NoError(t, c.Validate(req)) {
		if assert.NoError(t, routes.Create(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	}

	mockService.AssertExpectations(t)
}

func TestGetMessageByID(t *testing.T) {
	e, mockService, routes := setup()

	id := uuid.New()
	expectedMessage := entity.Message{ID: id, Message: "Hello, world!"}
	mockService.On("GetMessageById", mock.Anything, id).Return(expectedMessage, nil)

	req := httptest.NewRequest(http.MethodGet, "/messages/"+id.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	if assert.NoError(t, routes.GetByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var message entity.Message
		if assert.NoError(t, json.NewDecoder(rec.Body).Decode(&message)) {
			assert.Equal(t, expectedMessage, message)
		}
	}

	mockService.AssertExpectations(t)
}

func TestGetMessages(t *testing.T) {
	e, mockService, routes := setup()

	expectedMessages := []entity.Message{
		{ID: uuid.New(), Message: "Hello, world!"},
		{ID: uuid.New(), Message: "Hello, universe!"},
	}
	mockService.On("GetMessages", mock.Anything).Return(expectedMessages, nil)

	req := httptest.NewRequest(http.MethodGet, "/messages", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, routes.GetAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var messages []entity.Message
		if assert.NoError(t, json.NewDecoder(rec.Body).Decode(&messages)) {
			assert.Equal(t, expectedMessages, messages)
		}
	}

	mockService.AssertExpectations(t)
}

func TestMarkMessageAsProcessed(t *testing.T) {
	e, mockService, routes := setup()

	id := uuid.New()
	mockService.On("MarkMessageAsProcessed", mock.Anything, id).Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/messages/"+id.String()+"/process", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	if assert.NoError(t, routes.MarkAsProcessed(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	mockService.AssertExpectations(t)
}

func TestGetProcessedMessagesStats(t *testing.T) {
	e, mockService, routes := setup()

	expectedCount := 42
	mockService.On("GetProcessedMessagesStats", mock.Anything).Return(expectedCount, nil)

	req := httptest.NewRequest(http.MethodGet, "/messages/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, routes.GetStats(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response struct {
			ProcessedMessages int `json:"processed_messages"`
		}
		if assert.NoError(t, json.NewDecoder(rec.Body).Decode(&response)) {
			assert.Equal(t, expectedCount, response.ProcessedMessages)
		}
	}

	mockService.AssertExpectations(t)
}
