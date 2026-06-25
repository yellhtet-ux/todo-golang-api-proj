package todos

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/yellhtet-ux/todo-golang-api-proj/internal/adapters/postgresql/sqlc"
	"github.com/yellhtet-ux/todo-golang-api-proj/internal/json"
	"github.com/yellhtet-ux/todo-golang-api-proj/internal/todos/dto"
	todos "github.com/yellhtet-ux/todo-golang-api-proj/internal/todos/service"
)

type handler struct {
	service todos.Service
}

func NewHandler(service todos.Service) *handler {
	return &handler{
		service: service,
	}
}

// ListTodos godoc
// @Summary      List all todos
// @Description  Returns every todo in the system
// @Tags         todos
// @Produce      json
// @Success      200  {array}   dto.TodoResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /todos [get]
func (h *handler) ListTodos(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r,"userid")
	var userID pgtype.UUID

	if err := userID.Scan(id); err != nil {
		log.Println(err)
		json.InvalidRequest(w,err)
	}

	products, err := h.service.ListToDos(r.Context(),userID)

	if err != nil {
		log.Println(err)
		json.InternalServerError(w,err)
		return
	}
	json.Write(w, http.StatusOK, products)
}

// ListToDosByID godoc
// @Summary      Get todo by ID
// @Description  Returns a single todo matching the given UUID
// @Tags         todos
// @Produce      json
// @Param        id   path      string  true  "Todo UUID"
// @Success      200  {object}  dto.TodoResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /todo/{id} [get]
func (h *handler) ListToDosByID(w http.ResponseWriter, r *http.Request) {
	var param repo.ListToDosByIDParams
	var idParam = chi.URLParam(r, "id")
	var id pgtype.UUID
	if err := id.Scan(idParam); err != nil {
		log.Println(err)
		json.InvalidRequest(w, err)
		return
	}
		param = repo.ListToDosByIDParams {
			ID: id,
			UserID: id,
		}

	todo, err := h.service.ListToDosByID(r.Context(), param)
	if err != nil {
		log.Println(err)
		json.InternalServerError(w, err)
		return
	}
	json.Write(w, http.StatusOK, todo)
}

// CreateTodo godoc
// @Summary      Create a todo
// @Accept       json
// @Produce      json
// @Param        todo  body      dto.CreateTodoRequestDoc  true  "Todo to create"
// @Success      201   {string}  string  "Todo created successfully"
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /todo/create [post]
func (h *handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo *dto.CreateTodoRequest

	if err := json.Read(r, &todo); err != nil {
		log.Println(err)
		json.InvalidRequest(w, err)
		return
	}

	todoParams := &repo.CreateToDoParams{
		UserID: 		 todo.UserID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      repo.TodoStatus(todo.Status),
		Priority:    repo.TodoPriority(todo.Priority),
		DueAt:       pgtype.Timestamptz{Time: todo.DueDate, Valid: true},
	}

	log.Printf("To Do Params:%+v\n", todoParams)
	_, err := h.service.CreateTodo(r.Context(), todoParams)
	if err != nil {
		log.Println(err)
		json.InternalServerError(w, err)
		return
	}
	json.Write(w, http.StatusCreated, "Todo created successfully")
}

// UpdateTodoByStatus godoc
// @Summary      Update todo status
// @Description  Updates the status of a todo by ID
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id      path      string                true  "Todo UUID"
// @Param        status  body      dto.UpdateTodoStatus  true  "New status"
// @Success      200     {object}  dto.TodoResponse
// @Failure      400     {object}  dto.ErrorResponse
// @Failure      500     {object}  dto.ErrorResponse
// @Router       /todo/update/status/{id} [put]
func (h *handler) UpdateTodoByStatus(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	var id pgtype.UUID

	if err := id.Scan(idParam); err != nil {
		log.Println(err)
		json.InvalidRequest(w, err)
		return
	}

	var updatedStatus dto.UpdateTodoStatus

	if err := json.Read(r, &updatedStatus); err != nil {
		log.Println(err)
		json.InvalidRequest(w, err)
		return
	}

	updatedStatusParams := &repo.UpdateToDoStatusParams{
		ID:     id,
		Status: repo.TodoStatus(updatedStatus.Status),
	}

	updatedTodo, err := h.service.UpdateTodoByStatus(r.Context(), updatedStatusParams)
	if err != nil {
		log.Println(err)
		json.InternalServerError(w, err)
		return
	}

	json.Write(w, http.StatusOK, updatedTodo)
}

// UpdateToDoByPriority godoc
// @Summary      Update todo priority
// @Description  Updates the priority of a todo by ID
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id        path      string                  true  "Todo UUID"
// @Param        priority  body      dto.UpdateTodoPriority  true  "New priority"
// @Success      200       {object}  dto.TodoResponse
// @Failure      400       {object}  dto.ErrorResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /todo/update/priority/{id} [put]
func (h *handler) UpdateToDoByPriority(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	var id pgtype.UUID
	if err := id.Scan(idParam); err != nil {
		log.Println(err)
		json.InvalidRequest(w, err)
		return
	}

	var updatedPriority dto.UpdateTodoPriority

	if err := json.Read(r, &updatedPriority); err != nil {
		log.Println(err)
		json.InvalidRequest(w, err)
		return
	}

	updatedPriorityParam := &repo.UpdateToDoPriorityParams{
		ID:       id,
		Priority: repo.TodoPriority(updatedPriority.Priority),
	}

	todo, err := h.service.UpdateToDoByPriority(r.Context(), updatedPriorityParam)
	if err != nil {
		log.Println(err)
		json.InternalServerError(w, err)
		return
	}

	json.Write(w, http.StatusOK, todo)
}

// DeleteTodoByID godoc
// @Summary      Delete a todo
// @Description  Soft-deletes a todo by ID
// @Tags         todos
// @Produce      json
// @Param        id   path      string  true  "Todo UUID"
// @Success      200  {string}  string  "Todo deleted successfully"
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /todo/delete/{id} [delete]
func (h *handler) DeleteTodoByID(w http.ResponseWriter, r *http.Request) {
	var param repo.DeleteTodoByIDParams
	idParam := chi.URLParam(r, "id")
	var id pgtype.UUID
	if err := id.Scan(idParam); err != nil {
		log.Println(err)
		json.InvalidRequest(w, err)
		return
	}

	param = repo.DeleteTodoByIDParams{
		ID: id,
		UserID: id,
	}

	err := h.service.DeleteTodoByID(r.Context(), param)
	if err != nil {
		log.Println(err)
		json.InternalServerError(w, err)
		return
	}
	json.Write(w, http.StatusOK, "Todo deleted successfully")
}
