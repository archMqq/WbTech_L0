package order

import (
	"net/http"

	"github.com/gorilla/mux"
)

func GetOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	//TODO получение данных из КЭШа или БД
}
