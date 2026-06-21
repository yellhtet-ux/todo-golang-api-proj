package todos

import (
	"log"
	"net/http"

	// "time"

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
/// GET ENDPOINTS
// GET ALL TODOS
func (h *handler) ListTodos(w http.ResponseWriter,r *http.Request)  {
	products,err := h.service.ListToDos(r.Context())

	if err != nil {
		log.Println(err)
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	json.Write(w,http.StatusOK,products)
}

// GET TODO BY ID
func (h *handler) ListToDosByID(w http.ResponseWriter,r *http.Request)  {
	idParam := chi.URLParam(r,"id")
	var id pgtype.UUID
	if err := id.Scan(idParam); err != nil {
		log.Println(err)
		json.InvalidRequest(w,err)
		return
	}
	todo, err := h.service.ListToDosByID(r.Context(),id)
	if err != nil {
		log.Println(err)
		json.InternalServerError(w,err)
		return
	}
	json.Write(w,http.StatusOK,todo)
}

	// POST ENDPOINTS
	// CREATE TODO
	func (h *handler) CreateTodo(w http.ResponseWriter,r *http.Request)  {
		var todo dto.CreateTodoRequest

		if err := json.Read(r,&todo); err != nil {
			log.Println(err)
			json.InvalidRequest(w,err)
			return
		}

		todoParams := &repo.CreateToDoParams{
			Title: todo.Title,
			Description: todo.Description,
			Status: repo.TodoStatus(todo.Status),
			Priority: repo.TodoPriority(todo.Priority),
			DueAt: pgtype.Timestamptz{Time: todo.DueDate, Valid: true},
		}

		log.Printf("To Do Params:%+v\n", todoParams)
		_ , err := h.service.CreateTodo(r.Context(),todoParams)
		if err != nil {
			log.Println(err)
			json.InternalServerError(w,err)
			return
		}
		json.Write(w,http.StatusCreated,"Todo created successfully")
	}

	// UPDATE ENDPOINTS
	func (h *handler) UpdateTodo(w http.ResponseWriter,r *http.Request) {
		idParam := chi.URLParam(r, "id")
		var id pgtype.UUID

		if err := id.Scan(idParam); err != nil {
			log.Println(err)
			json.InvalidRequest(w,err)
			return 
		}

		var updatedStatus dto.UpdateTodoPriority

		if err := json.Read(r,&updatedStatus); err != nil {
			log.Println(err)
			json.InvalidRequest(w,err)
			return
		}

		updatedStatusParams := &repo.UpdateToDoStatusParams{
			ID: id,
			Status: repo.TodoStatus(updatedStatus.Status),
		}

		updatedTodo, err := h.service.UpdateTodoByID(r.Context(),updatedStatusParams)
		if err != nil {
			log.Println(err)
			json.InternalServerError(w,err)
		}

		json.Write(w,http.StatusOK,updatedTodo)
	}
	
	// DELETE ENDPOINTS
	// DELETE TODO BY ID
	func (h *handler) DeleteTodoByID(w http.ResponseWriter,r *http.Request) {
		idParam := chi.URLParam(r,"id")
		var id pgtype.UUID
		if err := id.Scan(idParam); err != nil {
			log.Println(err)
			json.InvalidRequest(w,err)
			return
		}

		err := h.service.DeleteTodoByID(r.Context(),id);
		if err != nil {
			log.Println(err)
			json.InternalServerError(w,err)
			return
		}
		json.Write(w,http.StatusOK,"Todo deleted successfully")
	}