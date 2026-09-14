package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example/input/abstract"
	"example/input/port"
)

type UserHandler struct {
	*abstract.AbstractHandler
	userUsecase port.UserUsecase
}

func NewUserHandler(useCase port.UserUsecase, oAbstractHandler *abstract.AbstractHandler) *UserHandler {
	return &UserHandler{
		AbstractHandler: oAbstractHandler,
		userUsecase:     useCase,
	}
}

func (oSelf *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	oQuery := r.URL.Query()
	name := oQuery.Get("name")

	user, err := oSelf.userUsecase.CreateUser(name)
	if err != nil {
		sErrorMessage := err.Error()
		http.Error(w, sErrorMessage, 500)
		return
	}

	oEncoder := json.NewEncoder(w)
	oEncoder.Encode(user)
}

func (oSelf *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")

	id, _ := strconv.Atoi(idStr)

	user, err := oSelf.userUsecase.GetUser(id)
	if err != nil {
		sErrorMessage := err.Error()
		http.Error(w, sErrorMessage, 404)
		return
	}

	oEncoder := json.NewEncoder(w)
	oEncoder.Encode(user)
}
