package main

import (
	"log"
	"net/http"

	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swaggo/http-swagger"

	repo "github.com/yellhtet-ux/todo-golang-api-proj/internal/adapters/postgresql/sqlc"
	"github.com/yellhtet-ux/todo-golang-api-proj/internal/todos/handler"
	todosService "github.com/yellhtet-ux/todo-golang-api-proj/internal/todos/service"
	"github.com/yellhtet-ux/todo-golang-api-proj/internal/user"

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

	r.Get("/v1/health", healthCheck)

	r.Get("/v1/swagger/*", httpSwagger.Handler())


	// User Related Routes
	userRepo := repo.New(app.db)
	userService := user.NewService(userRepo,app.db)
	userHandler := user.NewHandler(userService)

	r.Post("/v1/user/signup",userHandler.CreateUser)


	// Todos Related Routes
	todoRepo := repo.New(app.db)
	todoService := todosService.NewService(todoRepo)
	todoHandler := todos.NewHandler(todoService)

	// GET /todos
	r.Get("/v1/todos/{userid}", todoHandler.ListTodos)
	r.Get("/v1/todo/{id}", todoHandler.ListToDosByID)

	// POST /todo/create
	r.Post("/v1/todo/create", todoHandler.CreateTodo)

	// PUT /todo/update/status/{id}
	r.Put("/v1/todo/update/status/{id}", todoHandler.UpdateTodoByStatus)
	// PUT /todo/update/priority/{id}
	r.Put("/v1/todo/update/priority/{id}", todoHandler.UpdateToDoByPriority)

	// DELETE /todo/delete/{id}
	r.Delete("/v1/todo/delete/{id}", todoHandler.DeleteTodoByID)
	
	return r;
}

// healthCheck godoc
// @Summary      Health check
// @Description  Returns OK when the server is running
// @Tags         health
// @Produce      plain
// @Success      200  {string}  string  "ALL GOOD"
// @Router       /health [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ALL GOOD"))
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
