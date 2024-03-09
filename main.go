package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"hackowasp_ems/db"
	"log"
	"net/http"

	"hackowasp_ems/handler"

	"github.com/rs/cors"
)

func main() {
	db.InitDatabase()
	router := mux.NewRouter()
	router.HandleFunc("/team/{id}", handler.GetTeamByIDHandler).Methods("GET")
	router.HandleFunc("/teamAll", handler.GetAllTeamsHandler).Methods("GET")
	router.HandleFunc("/memberAll", handler.GetAllMembersHandler).Methods("GET")
	router.HandleFunc("/member/{id}", handler.GetMemberByIDHanlder).Methods("GET")
	router.HandleFunc("/member/{id}", handler.UpdateAttendanceHandler).Methods("PUT")

	corsMiddlerware := cors.Default()
	handler := corsMiddlerware.Handler(router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	go func() {
		fmt.Println("Server is running on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v\n", err)
		}
	}()
	WaitForTerminationSignal(server)

}

func WaitForTerminationSignal(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}
	log.Println("Server stopped gracefully")
}
