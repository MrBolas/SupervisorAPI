package handlers

import (
	"log"
	"net/http"

	"github.com/MrBolas/SupervisorAPI/auth"
	"github.com/MrBolas/SupervisorAPI/models"
	"github.com/MrBolas/SupervisorAPI/repositories"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/gofrs/uuid"
)

type TasksHandler struct {
	repo repositories.Repository
}

func NewTasksHandler(repo repositories.Repository) *TasksHandler {
	return &TasksHandler{
		repo: repo,
	}
}

func (th *TasksHandler) GetTaskById(c echo.Context) error {

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, "Invalid Id")
	}

	task, err := th.repo.GetTaskById(id)
	if err == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusNotFound, "Task not found")
	}
	if err != nil {
		return err
	}

	// If User does not have manager Role or owns task is unAuthorized
	if !auth.IsManager(c) && auth.GetUserId(c) != task.WorkerId {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	return c.JSON(http.StatusOK, task.ToResponse())
}

func (th *TasksHandler) GetTaskList(c echo.Context) error {

	query := repositories.NewListQuery()
	isManager := auth.IsManager(c)

	// auth
	if !isManager {
		query.Filters["worker_name"] = auth.GetUserNickname(c)
	}

	// pagination
	err := query.AddPageAndPageSize(c.QueryParam("page"), c.QueryParam("page_size"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// sorting
	err = query.AddSorting(c.QueryParam("sort_by"), c.QueryParam("sort_order"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// create filters
	err = query.AddListTaskFilters(c.QueryParams(), isManager)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Call to repository
	tasks, err := th.repo.ListTasks(query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.ToListResponse(tasks, query.Pagination.Page, query.Pagination.PageSize))
}

func (th *TasksHandler) CreateTask(c echo.Context) error {

	req := new(models.TaskRequest)
	err := c.Bind(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Malformed JSON")
	}

	err = req.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	task, err := req.ToTask(auth.GetUserNickname(c))
	if err != nil {
		return err
	}

	task, err = th.repo.CreateTask(task)
	if err == gorm.ErrRegistered {
		return c.JSON(http.StatusConflict, err)
	}
	if err != nil {
		return err
	}

	// Add task to Queue
	log.Println("The tech", task.WorkerId, "performed the task", task.Id.String(), "on date", task.Date.Time.String())

	return c.JSON(http.StatusCreated, task.ToResponse())
}

func (th *TasksHandler) DeleteTask(c echo.Context) error {

	// Only Manager can delete
	if !auth.IsManager(c) {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, "Invalid Id")
	}

	_, err = th.repo.GetTaskById(id)
	if err == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusNotFound, "Task not found")
	}
	if err != nil {
		return err
	}

	err = th.repo.DeleteTask(id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
