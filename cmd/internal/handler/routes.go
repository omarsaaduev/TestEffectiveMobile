package handler

import (
	_ "TestEffectiveMobile/docs"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(r *mux.Router, handler PersonHandler) {
	r.HandleFunc("/persons/", handler.GetPersons).Methods("GET")
	r.HandleFunc("/persons/", handler.AddPerson).Methods("POST")
	r.HandleFunc("/persons/{id}/", handler.UpdatePerson).Methods("PUT")
	r.HandleFunc("/persons/{id}/", handler.DeletePerson).Methods("DELETE")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
