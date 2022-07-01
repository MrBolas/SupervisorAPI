package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/MrBolas/SupervisorAPI/models"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

var mockedTask = models.Task{
	WorkerId: "mocked_worker_name",
	Summary:  "mocked_summary",
	Date: sql.NullTime{
		Time:  time.Date(2020, time.April, 21, 10, 10, 10, 0, time.UTC),
		Valid: true,
	},
}

var mockedTaskRequest = models.TaskRequest{
	Summary: "mocked_summary",
	Date:    "2006-01-02 03:04:05PM",
}

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "5.7",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=secret",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	databaseUrl := fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=True", resource.GetPort("3306/tcp"))

	log.Println("Connecting to test database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = gorm.Open(mysql.Open(databaseUrl), &gorm.Config{})
		if err != nil {
			return err
		}

		db.AutoMigrate(models.Task{})

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func teardown(t *testing.T) {
	sql, _ := db.DB()
	_, err := sql.Exec("DELETE FROM tasks")
	assert.Nil(t, err)
}

func TestCreateNewTask(t *testing.T) {

	mockedRepo := NewTasksRepository(db)
	defer teardown(t)

	createdTask, err := mockedRepo.CreateTask(mockedTask)
	assert.Nil(t, err)

	assert.Equal(t, createdTask.Summary, mockedTask.Summary)
	assert.Equal(t, createdTask.WorkerId, mockedTask.WorkerId)
	assert.Equal(t, createdTask.Date.Time, mockedTask.Date.Time)
	assert.Equal(t, createdTask.Date.Valid, mockedTask.Date.Valid)
}

func TestGetTaskById(t *testing.T) {

	mockedRepo := NewTasksRepository(db)
	defer teardown(t)

	createdTask, err := mockedRepo.CreateTask(mockedTask)
	assert.Nil(t, err)

	fetchedTenant, err := mockedRepo.GetTaskById(createdTask.Id)
	assert.Nil(t, err)

	assert.Equal(t, createdTask.Id, fetchedTenant.Id)
	assert.Equal(t, createdTask.Summary, fetchedTenant.Summary)
	assert.Equal(t, createdTask.WorkerId, fetchedTenant.WorkerId)
	assert.Equal(t, createdTask.Date.Time, fetchedTenant.Date.Time)
	assert.Equal(t, createdTask.Date.Valid, fetchedTenant.Date.Valid)
}

func TestGetTaskList(t *testing.T) {

	mockedRepo := NewTasksRepository(db)
	defer teardown(t)

	for i := 0; i < 10; i++ {
		newMockedTaskRequest, err := mockedTaskRequest.ToTask("mocked_worker_name")
		assert.Nil(t, err)

		_, err = mockedRepo.CreateTask(newMockedTaskRequest)
		assert.Nil(t, err)
	}

	query := NewListQuery()
	query.AddPageAndPageSize("", "")
	query.AddSorting("", "")

	tasks, err := mockedRepo.ListTasks(query)
	assert.Nil(t, err)

	assert.Equal(t, len(tasks), 10)
}

func TestGetTaskListByWorkerName(t *testing.T) {

	mockedRepo := NewTasksRepository(db)
	defer teardown(t)

	for i := 0; i < 5; i++ {
		newMockedTaskRequest, err := mockedTaskRequest.ToTask("mocked_worker_name_1")
		assert.Nil(t, err)

		_, err = mockedRepo.CreateTask(newMockedTaskRequest)
		assert.Nil(t, err)
	}

	for i := 0; i < 5; i++ {
		newMockedTaskRequest, err := mockedTaskRequest.ToTask("mocked_worker_name_2")
		assert.Nil(t, err)

		_, err = mockedRepo.CreateTask(newMockedTaskRequest)
		assert.Nil(t, err)
	}

	query := NewListQuery()
	query.AddPageAndPageSize("", "")
	query.AddSorting("", "")
	query.Filters["worker_name"] = "mocked_worker_name_1"

	tasks, err := mockedRepo.ListTasks(query)
	assert.Nil(t, err)

	assert.Equal(t, len(tasks), 5)
}

func TestGetTaskListBeforeTimestamp(t *testing.T) {

	mockedRepo := NewTasksRepository(db)
	defer teardown(t)

	for i := 0; i < 5; i++ {
		newMockedTaskRequest, err := mockedTaskRequest.ToTask("mocked_worker_name_1")
		assert.Nil(t, err)

		newMockedTaskRequest.Date.Time = time.Date(2020, time.April, 15, 10, 50, 0, 0, time.UTC)
		_, err = mockedRepo.CreateTask(newMockedTaskRequest)
		assert.Nil(t, err)
	}

	for i := 0; i < 5; i++ {
		newMockedTaskRequest, err := mockedTaskRequest.ToTask("mocked_worker_name_2")
		assert.Nil(t, err)

		_, err = mockedRepo.CreateTask(newMockedTaskRequest)
		assert.Nil(t, err)
	}

	query := NewListQuery()
	query.AddPageAndPageSize("", "")
	query.AddSorting("", "")

	urlValues := make(map[string][]string, 0)
	urlValues["before"] = []string{"2015-01-02 04:04:05PM"}
	query.AddListTaskFilters(urlValues, false)

	tasks, err := mockedRepo.ListTasks(query)
	assert.Nil(t, err)

	assert.Equal(t, len(tasks), 5)
	assert.Equal(t, tasks[0].WorkerId, "mocked_worker_name_2")
}

func TestGetTaskListAfterTimestamp(t *testing.T) {

	mockedRepo := NewTasksRepository(db)
	defer teardown(t)

	for i := 0; i < 5; i++ {
		newMockedTaskRequest, err := mockedTaskRequest.ToTask("mocked_worker_name_1")
		assert.Nil(t, err)

		newMockedTaskRequest.Date.Time = time.Date(2020, time.April, 15, 10, 50, 0, 0, time.UTC)
		_, err = mockedRepo.CreateTask(newMockedTaskRequest)
		assert.Nil(t, err)
	}

	for i := 0; i < 5; i++ {
		newMockedTaskRequest, err := mockedTaskRequest.ToTask("mocked_worker_name_2")
		assert.Nil(t, err)

		_, err = mockedRepo.CreateTask(newMockedTaskRequest)
		assert.Nil(t, err)
	}

	query := NewListQuery()
	err := query.AddPageAndPageSize("", "")
	assert.Nil(t, err)
	err = query.AddSorting("", "")
	assert.Nil(t, err)

	urlValues := make(map[string][]string, 0)
	urlValues["after"] = []string{"2015-01-02 04:04:05PM"}
	err = query.AddListTaskFilters(urlValues, false)
	assert.Nil(t, err)

	tasks, err := mockedRepo.ListTasks(query)
	assert.Nil(t, err)

	assert.Equal(t, len(tasks), 5)
	assert.Equal(t, tasks[0].WorkerId, "mocked_worker_name_1")
}

func TestUpdateTask(t *testing.T) {

	mockedRepo := NewTasksRepository(db)
	defer teardown(t)

	newTask, err := mockedTaskRequest.ToTask("mocked_worker_name")
	assert.Nil(t, err)

	createdTask, err := mockedRepo.CreateTask(newTask)
	assert.Nil(t, err)

	modifiedTask := createdTask
	modifiedTask.Summary = "updated mocked summary text"

	updatedTask, err := mockedRepo.UpdateTask(createdTask.Id, createdTask, modifiedTask)
	assert.Nil(t, err)

	assert.Equal(t, modifiedTask.Id, updatedTask.Id)
	assert.Equal(t, modifiedTask.Summary, updatedTask.Summary)
	assert.Equal(t, modifiedTask.WorkerId, updatedTask.WorkerId)
	assert.Equal(t, modifiedTask.Date.Time, updatedTask.Date.Time)
	assert.Equal(t, modifiedTask.Date.Valid, updatedTask.Date.Valid)
}
