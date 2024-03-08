package main

import (
	"fmt"

	"hackowasp_ems/db"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db.InitDatabase()
	router := mux.NewRouter()
	router.HandleFunc("/team/{id}", GetTeamByIDHandler).Methods("GET")
	router.HandleFunc("/member/{id}", GetMemberByIDHanlder).Methods("GET")
	router.HandleFunc("/member/{id}", UpdateAttendanceHandler).Methods("PUT")
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		fmt.Println("Server is running on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v\n", err)
		}
	}()
	waitForTerminationSignal(server)

}
