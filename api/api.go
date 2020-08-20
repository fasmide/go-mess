package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fasmide/go-mess/database"
)

type API struct {
	http.ServeMux
	Db *database.Database
}

func New(d *database.Database) *API {
	a := &API{
		Db: d,
	}

	a.HandleFunc("/active", a.active)
	a.HandleFunc("/previous", a.previous)
	return a
}

func (a *API) active(w http.ResponseWriter, r *http.Request) {
	data, err := a.Db.ActiveOrders()
	if err != nil {
		err = fmt.Errorf("could not lookup active orders: %w", err)
		log.Printf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(data)
	if err != nil {
		log.Printf("could not json encode active orders: %s", err)
	}
}

func (a *API) previous(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	ids, exists := values["id"]
	if !exists {
		http.Error(w, "no id query parameters provided", http.StatusBadRequest)
		return
	}

	data, err := a.Db.PreviousOrders(ids...)
	if err != nil {
		errStr := fmt.Sprintf("could not lookup previous orders: %s", err)
		log.Printf(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(data)
	if err != nil {
		log.Printf("could not encode previous orders: %s", err)
	}
}
