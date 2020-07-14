package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type SuggestionReq struct {
	Id     string `json:"id"`
	Email  string `json:"email"`
	Detail string `json:"detail"`
	Date   string `json:"date"`
}

type SuggestionRes struct {
	Id     string `json:"id"`
	Email  string `json:"email"`
	Detail string `json:"detail"`
	Date   string `json:"date"`
}

type Error struct {
	Message string `json:"message"`
}

var db *sql.DB

// main func

func main() {

	// DB driver folder
	pgUrl, err := pq.ParseURL("postgres:")

	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("postgres", pgUrl)

	if err != nil {
		log.Fatal(err)
	}

	// router / handler definition - folder
	router := mux.NewRouter()
	router.HandleFunc("/v1/Feedback/InitialFeedback", insertHandler).Methods("POST")

	//port setup
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting port to %s", port)
	}
	log.Printf("Listening to port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}

}

// insertHandler func

func insertHandler(w http.ResponseWriter, r *http.Request) {

	// controller - folder
	fmt.Println("In the IH")

	var iHSuggestionReq SuggestionReq
	var iHSuggestionRes SuggestionRes
	var error Error

	fmt.Println("iHSuggestionReq:", iHSuggestionReq)
	fmt.Println("r.Body:", r.Body)

	err := json.NewDecoder(r.Body).Decode(&iHSuggestionReq)

	if err != nil {
		error.Message = "Bad data"
		responseWithError(w, http.StatusBadRequest, error)
		return
	}

	if iHSuggestionReq.Email == "" {
		error.Message = "Email ID should not be empty"
		responseWithError(w, http.StatusBadRequest, error)
		return
	}

	fmt.Println("iHSugestiontReq:", iHSuggestionReq)

	queryDet := "insert into userfeedback (email, detail, date) values($1, $2, $3) RETURNING id;"

	err1 := db.QueryRow(queryDet, iHSuggestionReq.Email, iHSuggestionReq.Detail, iHSuggestionReq.Date).Scan(&iHSuggestionRes.Id)

	if err1 != nil {
		log.Fatal(err1)
	}

	iHSuggestionRes.Email = iHSuggestionReq.Email
	iHSuggestionRes.Detail = iHSuggestionReq.Detail
	iHSuggestionRes.Date = iHSuggestionReq.Date
	//	iHSuggestionRes.Id = "1"

	w.Header().Set("content-type", "application/json")

	json.NewEncoder(w).Encode(iHSuggestionRes)

}

// responseWithError func
func responseWithError(w http.ResponseWriter, status int, error Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}
