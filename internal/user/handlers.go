package user

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yellhtet-ux/todo-golang-api-proj/internal/json"
)

type handler struct {
	service Service
}

func NewHandler (service Service) *handler {
	return &handler {
		service: service,
	}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request)  {
	var params createUserParam 

	if err := json.Read(r,&params); err != nil {
		log.Println(err)
		json.InvalidRequest(w,err)		
	}

	user,err := h.service.CreateUser(r.Context(),params)

	if err != nil {
		log.Println(err)
		json.InternalServerError(w,err)
	}

	json.Write(w,http.StatusCreated,user);
}

func (h *handler) GetUserByID(w http.ResponseWriter,r *http.Request) {	
	id := chi.URLParam(r,"userid")
	var userID pgtype.UUID

	if err := userID.Scan(id); err != nil {
			log.Println(err)
			json.InvalidRequest(w,err)
	}

	user, err := h.service.GetUserByID(r.Context(),userID)
	
	if err != nil {
		log.Println(err)
		json.InternalServerError(w,err)
	}

	json.Write(w,http.StatusOK,user)
}
