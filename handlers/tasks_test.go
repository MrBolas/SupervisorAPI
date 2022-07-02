package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MrBolas/SupervisorAPI/encryption"
	"github.com/MrBolas/SupervisorAPI/models"
	"github.com/MrBolas/SupervisorAPI/repositories"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockRepo struct {
	mock.Mock
}

var (
	mockedTaskRequest = models.TaskRequest{
		Summary: "wTThqMkifM_XNUE8WPnFLjhDOIlGD9cur5loFiQN",
		Date:    "2022-05-23 03:33:01PM",
	}
	mockedTask = models.Task{
		Id:       uuid.FromStringOrNil("a2d45497-09b4-4da1-a0d0-173d0bd12f13"),
		WorkerId: "mocked_worker_id",
		Summary:  "wTThqMkifM_XNUE8WPnFLjhDOIlGD9cur5loFiQN",
		Date: sql.NullTime{
			Valid: true,
			Time:  time.Date(2020, time.April, 15, 10, 50, 0, 0, time.UTC),
		},
	}
)

func (mr *mockRepo) GetTaskById(id uuid.UUID) (models.Task, error) {
	args := mr.Called(id)

	mockedID := args.Get(0)
	if mockedID == nil {
		return models.Task{}, args.Error(1)
	}

	return args.Get(0).(models.Task), args.Error(1)
}

func (mr *mockRepo) CreateTask(t models.Task) (models.Task, error) {
	args := mr.Called(t)

	mockedTask := args.Get(0)
	if mockedTask == nil {
		return models.Task{}, args.Error(1)
	}

	return args.Get(0).(models.Task), args.Error(1)
}

func (mr *mockRepo) ListTasks(filters repositories.ListQuery) ([]models.Task, error) {
	args := mr.Called()

	mockedTask := args.Get(0)
	if mockedTask == nil {
		return []models.Task{}, args.Error(1)
	}

	return args.Get(0).([]models.Task), args.Error(1)
}

func (mr *mockRepo) UpdateTask(id uuid.UUID, oldTask models.Task, newTask models.Task) (models.Task, error) {
	args := mr.Called(id)

	mockedID := args.Get(0)
	if mockedID == nil {
		return models.Task{}, args.Error(1)
	}

	return args.Get(0).(models.Task), args.Error(1)
}

func (mr *mockRepo) DeleteTask(id uuid.UUID) error {
	args := mr.Called(id)

	return args.Error(0)
}

func addClaimsToJWTContext(c echo.Context, mockedClaims map[string]string) {

	// Add role to claim
	claims := make(jwt.MapClaims)

	for mockedKey, mockedValue := range mockedClaims {
		claims[mockedKey] = mockedValue
	}

	// define jwt token
	tk := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// add jwt to context
	c.Set("user", tk)
}

func createRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6357",
		Password: "",
		DB:       0,
	})
}

func TestGetTaskByIdShould200OK(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/tasks/a2d45497-09b4-4da1-a0d0-173d0bd12f13", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := make(map[string]string, 0)
	claims["http://supervisorapi/role"] = "manager"
	addClaimsToJWTContext(c, claims)

	c.SetPath("/tasks/:id")
	c.SetParamNames("id")
	c.SetParamValues("a2d45497-09b4-4da1-a0d0-173d0bd12f13")

	mr := mockRepo{}
	ce := encryption.NewCryptoEngine("Qp7LtWv8X4xEHk8OLidUOCUHURPaBmPk")
	mr.On("GetTaskById", mock.Anything).Return(mockedTask, nil)
	h := NewTasksHandler(&mr, ce, createRedisClient())

	taskResponse := mockedTask.ToResponse()
	taskResponse.Summary = ce.Decrypt(taskResponse.Summary)
	u, err := json.Marshal(taskResponse)
	assert.Nil(t, err)

	// Assertions
	if assert.NoError(t, h.GetTaskById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(u)+"\n", rec.Body.String())
	}
}

func TestGetTaskByIdShould404NotFoundWhenDoesntExist(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/tasks/a2d45497-09b4-4da1-a0d0-173d0bd12f13", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := make(map[string]string, 0)
	claims["http://supervisorapi/role"] = "manager"
	addClaimsToJWTContext(c, claims)

	c.SetPath("/tasks/:id")
	c.SetParamNames("id")
	c.SetParamValues("a2d45497-09b4-4da1-a0d0-173d0bd12f13")

	mr := mockRepo{}
	ce := encryption.NewCryptoEngine("Qp7LtWv8X4xEHk8OLidUOCUHURPaBmPk")
	mr.On("GetTaskById", mock.Anything).Return(models.Task{}, gorm.ErrRecordNotFound)
	h := NewTasksHandler(&mr, ce, createRedisClient())

	// Assertions
	if assert.NoError(t, h.GetTaskById(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func TestCreateTaskShould201Created(t *testing.T) {
	e := echo.New()
	u, err := json.Marshal(mockedTaskRequest)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/task", strings.NewReader(string(u)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := make(map[string]string, 0)
	claims["http://supervisorapi/role"] = "manager"
	claims["http://supervisorapi/nickname"] = "mocked_worker_id"
	addClaimsToJWTContext(c, claims)

	c.SetPath("/tasks")

	mr := mockRepo{}
	ce := encryption.NewCryptoEngine("Qp7LtWv8X4xEHk8OLidUOCUHURPaBmPk")
	mr.On("CreateTask", mock.Anything).Return(mockedTask, nil)
	h := NewTasksHandler(&mr, ce, createRedisClient())

	mockedTaskResponse := mockedTask.ToResponse()
	u, err = json.Marshal(mockedTaskResponse)
	assert.Nil(t, err)

	// Assertions
	if assert.NoError(t, h.CreateTask(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, string(u)+"\n", rec.Body.String())
	}
}

func TestCreateTaskShould400BadRequestWhenTimeFormatIsInvalid(t *testing.T) {
	e := echo.New()
	modifiedmockedTaskRequest := mockedTaskRequest
	modifiedmockedTaskRequest.Date = "2022 09 12 02:40:30"
	u, err := json.Marshal(modifiedmockedTaskRequest)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/task", strings.NewReader(string(u)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := make(map[string]string, 0)
	claims["http://supervisorapi/role"] = "manager"
	claims["http://supervisorapi/nickname"] = "mocked_worker_id"
	addClaimsToJWTContext(c, claims)

	c.SetPath("/tasks")

	mr := mockRepo{}
	ce := encryption.NewCryptoEngine("Qp7LtWv8X4xEHk8OLidUOCUHURPaBmPk")
	mr.On("CreateTask", mock.Anything).Return(mockedTask, nil)
	h := NewTasksHandler(&mr, ce, createRedisClient())

	mockedTaskResponse := mockedTask.ToResponse()
	u, err = json.Marshal(mockedTaskResponse)
	assert.Nil(t, err)

	// Assertions
	if assert.NoError(t, h.CreateTask(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "\"invalid date format, use yyyy-mm-dd hh:mm:ssPM\"\n", rec.Body.String())
	}
}

func TestCreateTaskShould400BadRequestWhenSummaryIsInvalid(t *testing.T) {
	e := echo.New()
	modifiedmockedTaskRequest := mockedTaskRequest
	modifiedmockedTaskRequest.Summary = string(make([]byte, 2501))
	u, err := json.Marshal(modifiedmockedTaskRequest)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/task", strings.NewReader(string(u)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := make(map[string]string, 0)
	claims["http://supervisorapi/role"] = "manager"
	claims["http://supervisorapi/nickname"] = "mocked_worker_id"
	addClaimsToJWTContext(c, claims)

	c.SetPath("/tasks")

	mr := mockRepo{}
	ce := encryption.NewCryptoEngine("Qp7LtWv8X4xEHk8OLidUOCUHURPaBmPk")
	mr.On("CreateTask", mock.Anything).Return(mockedTask, nil)
	h := NewTasksHandler(&mr, ce, createRedisClient())

	mockedTaskResponse := mockedTask.ToResponse()
	u, err = json.Marshal(mockedTaskResponse)
	assert.Nil(t, err)

	// Assertions
	if assert.NoError(t, h.CreateTask(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "\"summary max size is 2500 characters\"\n", rec.Body.String())
	}
}

func TestCreateTaskShould409ConflictWhenTaskAlreadyExists(t *testing.T) {
	e := echo.New()

	u, err := json.Marshal(mockedTaskRequest)
	assert.Nil(t, err)
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(string(u)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := make(map[string]string, 0)
	claims["http://supervisorapi/role"] = "manager"
	claims["http://supervisorapi/nickname"] = "mocked_worker_id"
	addClaimsToJWTContext(c, claims)

	c.SetPath("/tasks")

	mr := mockRepo{}
	ce := encryption.NewCryptoEngine("Qp7LtWv8X4xEHk8OLidUOCUHURPaBmPk")
	mr.On("CreateTask", mock.Anything).Return(models.Task{}, gorm.ErrRegistered)
	h := NewTasksHandler(&mr, ce, createRedisClient())

	// Assertions
	if assert.NoError(t, h.CreateTask(c)) {
		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Equal(t, "\"conflict creating task\"\n", rec.Body.String())
	}
}

func TestGetTaskListShould200OKListingFilteredTasks(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := make(map[string]string, 0)
	claims["http://supervisorapi/role"] = "manager"
	claims["http://supervisorapi/nickname"] = "mocked_worker_id"
	addClaimsToJWTContext(c, claims)

	c.SetPath("/tasks")

	taskList := []models.Task{mockedTask, mockedTask, mockedTask}

	mr := mockRepo{}
	ce := encryption.NewCryptoEngine("Qp7LtWv8X4xEHk8OLidUOCUHURPaBmPk")
	mr.On("ListTasks", mock.Anything).Return(taskList, nil)
	h := NewTasksHandler(&mr, ce, createRedisClient())

	decryptedTaskList := []models.Task{}
	for _, task := range taskList {
		task.Summary = ce.Decrypt(task.Summary)
		decryptedTaskList = append(decryptedTaskList, task)
	}

	taskListResponse := models.ToListResponse(decryptedTaskList, 1, 20)
	u, err := json.Marshal(taskListResponse)
	assert.Nil(t, err)

	// Assertions
	if assert.NoError(t, h.GetTaskList(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(u)+"\n", rec.Body.String())
	}
}

func TestGetTaskListShould400BadRequestWhenPageNumberIsLessThan1(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/tenants", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := make(map[string]string, 0)
	claims["http://supervisorapi/role"] = "manager"
	claims["http://supervisorapi/nickname"] = "mocked_worker_id"
	addClaimsToJWTContext(c, claims)

	c.SetPath("/tasks")
	c.QueryParams().Add("page", "0")

	taskList := []models.Task{mockedTask, mockedTask, mockedTask}

	mr := mockRepo{}
	ce := encryption.NewCryptoEngine("Qp7LtWv8X4xEHk8OLidUOCUHURPaBmPk")
	mr.On("ListTasks", mock.Anything).Return(taskList, nil)
	h := NewTasksHandler(&mr, ce, createRedisClient())

	// Assertions
	if assert.NoError(t, h.GetTaskList(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "\"page must be bigger than 0\"\n", rec.Body.String())
	}
}

func TestGetTaskListShould400BadRequestWhenPageSizeNumberIsLessThan1(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/tenants", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := make(map[string]string, 0)
	claims["http://supervisorapi/role"] = "manager"
	claims["http://supervisorapi/nickname"] = "mocked_worker_id"
	addClaimsToJWTContext(c, claims)

	c.SetPath("/tasks")
	c.QueryParams().Add("page", "1")
	c.QueryParams().Add("page_size", "0")

	taskList := []models.Task{mockedTask, mockedTask, mockedTask}

	mr := mockRepo{}
	ce := encryption.NewCryptoEngine("Qp7LtWv8X4xEHk8OLidUOCUHURPaBmPk")
	mr.On("ListTasks", mock.Anything).Return(taskList, nil)
	h := NewTasksHandler(&mr, ce, createRedisClient())

	// Assertions
	if assert.NoError(t, h.GetTaskList(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "\"page_size must be bigger than 0\"\n", rec.Body.String())
	}
}

func TestGetTaskListShould400BadRequestWhenPageSizeNumberIsMoreThan40(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/tenants", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := make(map[string]string, 0)
	claims["http://supervisorapi/role"] = "manager"
	claims["http://supervisorapi/nickname"] = "mocked_worker_id"
	addClaimsToJWTContext(c, claims)

	c.SetPath("/tasks")
	c.QueryParams().Add("page", "1")
	c.QueryParams().Add("page_size", "41")

	taskList := []models.Task{mockedTask, mockedTask, mockedTask}

	mr := mockRepo{}
	ce := encryption.NewCryptoEngine("Qp7LtWv8X4xEHk8OLidUOCUHURPaBmPk")
	mr.On("ListTasks", mock.Anything).Return(taskList, nil)
	h := NewTasksHandler(&mr, ce, createRedisClient())

	// Assertions
	if assert.NoError(t, h.GetTaskList(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "\"page_size must be less or equal than 40\"\n", rec.Body.String())
	}
}

func TestUpdateTaskShould200OK(t *testing.T) {
	e := echo.New()
	u, err := json.Marshal(mockedTaskRequest)
	assert.Nil(t, err)
	req := httptest.NewRequest(http.MethodPut, "/tasks/a2d45497-09b4-4da1-a0d0-173d0bd12f13", strings.NewReader(string(u)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := make(map[string]string, 0)
	claims["http://supervisorapi/role"] = "manager"
	claims["http://supervisorapi/nickname"] = "mocked_worker_id"
	addClaimsToJWTContext(c, claims)

	c.SetPath("/tasks/:id")
	c.SetParamNames("id")
	c.SetParamValues("a2d45497-09b4-4da1-a0d0-173d0bd12f13")

	mr := mockRepo{}
	ce := encryption.NewCryptoEngine("Qp7LtWv8X4xEHk8OLidUOCUHURPaBmPk")

	updatedMockedTask := mockedTask
	updatedMockedTask.Summary = ce.Encrypt("updated mock summary")

	mr.On("GetTaskById", mock.Anything).Return(mockedTask, nil)
	mr.On("UpdateTask", updatedMockedTask.Id, mock.Anything).Return(updatedMockedTask, nil)
	h := NewTasksHandler(&mr, ce, createRedisClient())

	u, err = json.Marshal(updatedMockedTask.ToResponse())
	assert.Nil(t, err)

	// Assertions
	if assert.NoError(t, h.UpdateTask(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(u)+"\n", rec.Body.String())
	}
}

func TestUpdateTaskShould404NotFoundWhenTaskDoesNotExist(t *testing.T) {
	e := echo.New()
	u, err := json.Marshal(mockedTaskRequest)
	assert.Nil(t, err)
	req := httptest.NewRequest(http.MethodPut, "/tasks/a2d45497-09b4-4da1-a0d0-173d0bd12f13", strings.NewReader(string(u)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := make(map[string]string, 0)
	claims["http://supervisorapi/role"] = "manager"
	claims["http://supervisorapi/nickname"] = "mocked_worker_id"
	addClaimsToJWTContext(c, claims)

	c.SetPath("/tasks/:id")
	c.SetParamNames("id")
	c.SetParamValues("a2d45497-09b4-4da1-a0d0-173d0bd12f13")

	mr := mockRepo{}
	ce := encryption.NewCryptoEngine("Qp7LtWv8X4xEHk8OLidUOCUHURPaBmPk")
	mr.On("GetTaskById", mock.Anything).Return(models.Task{}, gorm.ErrRecordNotFound)
	h := NewTasksHandler(&mr, ce, createRedisClient())

	// Assertions
	if assert.NoError(t, h.UpdateTask(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "\"Not found\"\n", rec.Body.String())
	}
}

func TestDeleteTaskShould204NoContent(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/tasks/a2d45497-09b4-4da1-a0d0-173d0bd12f13", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := make(map[string]string, 0)
	claims["http://supervisorapi/role"] = "manager"
	claims["http://supervisorapi/nickname"] = "mocked_worker_id"
	addClaimsToJWTContext(c, claims)

	c.SetPath("/tasks/:id")
	c.SetParamNames("id")
	c.SetParamValues("a2d45497-09b4-4da1-a0d0-173d0bd12f13")

	mr := mockRepo{}
	mr.On("GetTaskById", mock.Anything).Return(mockedTask, nil)
	mr.On("DeleteTask", mock.Anything).Return(nil)
	ce := encryption.NewCryptoEngine("Qp7LtWv8X4xEHk8OLidUOCUHURPaBmPk")
	h := NewTasksHandler(&mr, ce, createRedisClient())

	// Assertions
	if assert.NoError(t, h.DeleteTask(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}

func TestDeleteTenantShouldReturn404NotFoundWhenItDoesntExist(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/tasks/a2d45497-09b4-4da1-a0d0-173d0bd12f13", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := make(map[string]string, 0)
	claims["http://supervisorapi/role"] = "manager"
	claims["http://supervisorapi/nickname"] = "mocked_worker_id"
	addClaimsToJWTContext(c, claims)

	c.SetPath("/tasks/:id")
	c.SetParamNames("id")
	c.SetParamValues("a2d45497-09b4-4da1-a0d0-173d0bd12f13")

	mr := mockRepo{}
	mr.On("GetTaskById", mock.Anything).Return(models.Task{}, gorm.ErrRecordNotFound)
	mr.On("DeleteTask", mock.Anything).Return(nil)
	ce := encryption.NewCryptoEngine("Qp7LtWv8X4xEHk8OLidUOCUHURPaBmPk")
	h := NewTasksHandler(&mr, ce, createRedisClient())

	// Assertions
	if assert.NoError(t, h.DeleteTask(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "\"Task not found\"\n", rec.Body.String())
	}
}
