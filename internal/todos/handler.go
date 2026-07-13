package todos

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yellhtet-ux/todo-golang-api-proj/env"
	"github.com/yellhtet-ux/todo-golang-api-proj/internal/json"
	"github.com/yellhtet-ux/todo-golang-api-proj/internal/user"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
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
// @Router       /v1/todos/{user_id} [get]
func (h *handler) ListTodos(w http.ResponseWriter, r *http.Request) {

	claimKey := env.GetString("USER_CLAIMS","user_claims")

	claims, ok := r.Context().Value(claimKey).(*user.CustomClaims)

	if !ok {
		json.InternalServerError(w,fmt.Errorf("failed to get user claims"),nil)
		return
	}

	todos, err := h.service.ListToDos(r.Context(),claims.UserID)

	if err != nil {
		log.Println(err)
		json.InternalServerError(w,err,nil)
		return
	}

	json.Write(w, http.StatusOK, todos)
}

// ListToDosByID godoc
// @Summary      Get todo by ID
// @Description  Returns a single todo matching the given UUID
// @Tags         todos
// @Produce      json
// @Param        id   path      string  true  "Todo UUID"
// @Param        user_id   path  string  true  "User UUID"
// @Success      200  {object}  dto.TodoResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /v1/todo/ [get]
func (h *handler) ListToDosByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r,"todo_id")
	var id pgtype.UUID

	if err := id.Scan(idParam); err != nil {
		log.Println(err)
		json.InvalidRequest(w,err,nil)
		return
	}

	claimKey := env.GetString("USER_CLAIMS","user_claims")

	claims , ok := r.Context().Value(claimKey).(*user.CustomClaims)

	if !ok {
		json.InternalServerError(w,fmt.Errorf("failed to get user claims"),nil)
		return 
	}

	todo, err := h.service.ListToDosByID(r.Context(), id, claims.UserID)

	if err != nil {
		log.Println(err)
		json.NotFound(w,err,nil)
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
	var todo CreateTodoRequest

	if err := json.Read(r, &todo); err != nil {
		log.Println(err)
		json.InvalidRequest(w, err,nil)
		return
	}

	claimKey := env.GetString("USER_CLAIMS","user_claims")

	claims, ok := r.Context().Value(claimKey).(*user.CustomClaims)

	if !ok {
		json.InternalServerError(w,fmt.Errorf("failed to get user claims"),nil)
		return
	}

 	_, err := h.service.CreateTodo(r.Context(),claims.UserID,todo)

	if err != nil {
		log.Println(err)
		json.NotFound(w,err,nil)
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
	idParam := chi.URLParam(r, "todo_id")
	var id pgtype.UUID

	if err := id.Scan(idParam); err != nil {
		log.Println(err)
		json.InvalidRequest(w, err,nil)
		return
	}

	claimKey := env.GetString("USER_CLAIMS","user_claims")

	claims, ok := r.Context().Value(claimKey).(*user.CustomClaims)

	if !ok {
		json.InternalServerError(w,fmt.Errorf("failed to get user claims"),nil)
		return
	}

	var updatedStatus UpdateToDoStatusRequest

	if err := json.Read(r, &updatedStatus); err != nil {
		log.Println(err)
		json.InvalidRequest(w, err,nil)
		return
	}

	updatedStatusParams := UpdateToDoStatusRequest{
			Status: updatedStatus.Status,
	} 

	updatedTodo, err := h.service.UpdateTodoByStatus(r.Context(),id,claims.UserID,updatedStatusParams)
	if err != nil {
		log.Println(err)
		json.NotFound(w, err,nil)
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
	idParam := chi.URLParam(r, "todo_id")
	var id pgtype.UUID
	if err := id.Scan(idParam); err != nil {
		log.Println(err)
		json.InvalidRequest(w, err,nil)
		return
	}

	claimKey := env.GetString("USER_CLAIMS","user_claims")

	claims, ok := r.Context().Value(claimKey).(*user.CustomClaims)

	if !ok {
		json.InternalServerError(w,fmt.Errorf("failed to get user claims"),nil)
		return
	}

	var updatedPriority UpdateToDoPriorityRequest 

	if err := json.Read(r, &updatedPriority); err != nil {
		log.Println(err)
		json.InvalidRequest(w, err,nil)
		return
	}

	updatedPriorityParam := UpdateToDoPriorityRequest {
		Priority: updatedPriority.Priority,
	}

	todo, err := h.service.UpdateToDoByPriority(r.Context(),id,claims.UserID,updatedPriorityParam)
	if err != nil {
		log.Println(err)
		json.NotFound(w, err,nil)
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
	idParam := chi.URLParam(r, "todo_id")
	var id pgtype.UUID
		if err := id.Scan(idParam); err != nil {
			log.Println(err)
			json.InvalidRequest(w, err,nil)
			return
		}

		claimKey := env.GetString("USER_CLAIMS","user_claims")
		
		claims, ok := r.Context().Value(claimKey).(*user.CustomClaims)

		if !ok {
			json.InternalServerError(w,fmt.Errorf("failed to get user claims"),nil)
			return
		}

		err := h.service.DeleteTodoByID(r.Context(),id,claims.UserID)

	if err != nil {
		log.Println(err)
		json.NotFound(w, err,nil)
		return
	}
	json.Write(w, http.StatusOK, "Todo deleted successfully")
}
