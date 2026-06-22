package main

import (
	"log"
	"net/http"

	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	repo "github.com/yellhtet-ux/todo-golang-api-proj/internal/adapters/postgresql/sqlc"
	"github.com/yellhtet-ux/todo-golang-api-proj/internal/todos/handler"
	todosService "github.com/yellhtet-ux/todo-golang-api-proj/internal/todos/service"
)


func (app *application) mount () http.Handler { 
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID) // Important for rate limiting
	r.Use(middleware.ClientIPFromRemoteAddr) // Important for rate limiting, analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // Recover from crashes

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ALL GOOD"))
	})

	todoRepo := repo.New(app.db)
	todoService := todosService.NewService(todoRepo)
	todoHandler := todos.NewHandler(todoService)

	// GET /todos
	r.Get("/todos", todoHandler.ListTodos)
	r.Get("/todo/{id}", todoHandler.ListToDosByID)

	// POST /todo/create
	r.Post("/todo/create", todoHandler.CreateTodo)

	// PUT /todo/update/{id}
	r.Put("/todo/update/{id}",todoHandler.UpdateTodoByStatus)
	r.Put("/todo/update/{id}",todoHandler.UpdateToDoByPriority)

	// DELETE /todo/delete/{id}
	r.Delete("/todo/delete/{id}", todoHandler.DeleteTodoByID)
	
	return r;
}


func (app *application) run (h http.Handler) error {
	srv := &http.Server {
		Addr: app.config.addr,
		Handler: h,
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout: time.Minute,
	}

	log.Printf("Server is running with this address %s",app.config.addr)

	return srv.ListenAndServe()
}

type application struct {
	config config
	db *pgxpool.Pool
}

type config struct {
	addr string
	db dbConfig
}

type dbConfig struct {
	dsn string
}