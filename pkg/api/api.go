package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func AddressClaimableHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "address :%v", vars["address"])
}

func NewApiServer() error {
	if err := godotenv.Load(".env"); err != nil {
		return err
	}

	fmt.Printf("Starting server at port : %s", os.Getenv("PORT"))

	r := mux.NewRouter()

	r.HandleFunc("/addresses/{address}/claimable", AddressClaimableHandler)

	http.Handle("/", r)

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		return err
	}
	return nil
}
