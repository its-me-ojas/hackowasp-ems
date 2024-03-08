package main

import (
	"context"
	"encoding/json"
	"fmt"

	"hackowasp_ems/db"
	"hackowasp_ems/model"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func UpdateAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	memberIDBytes := r.URL.Path[len("/member/"):]
	memberID := string(memberIDBytes)
	fmt.Println(memberID)

	if len(memberID) == 0 {
		http.Error(w, "Member ID not provided", http.StatusBadRequest)
		return
	}
	objID, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		http.Error(w, "Invalid Member ID", http.StatusBadRequest)
		return
	}

	// Retrieve the current member
	var member model.HackMember
	err = db.HackMemberCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&member)
	if err != nil {
		http.Error(w, "Error retrieving member: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error retrieving member:", err)
		return
	}

	// Toggle the checkIn field
	member.CheckedIn = !member.CheckedIn
	member.Time = time.Now()

	// Update the member in the database
	update := bson.M{
		"$set": member,
	}
	_, err = db.HackMemberCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		http.Error(w, "Error updating member: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error updating member:", err)
		return
	}

	fmt.Println("Member check-in status toggled successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Member check-in status toggled successfully"}`))
}

func GetTeamByIDHandler(w http.ResponseWriter, r *http.Request) {
	teamIDBytes := r.URL.Path[len("/team/"):]
	teamID := string(teamIDBytes)
	fmt.Println(teamID)
	if len(teamID) == 0 {
		http.Error(w, "Team ID not provided", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(teamID)
	if err != nil {
		http.Error(w, "Invalid Team ID", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objID}

	var team model.HackTeam
	err = db.HackTeamCollection.FindOne(context.Background(), filter).Decode(&team)
	if err != nil {
		http.Error(w, "Error retrieving team: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error retrieving team:", err)
		return
	}
	fmt.Println("Team Members:")
	for _, member := range team.Members {
		fmt.Println("Member:", member.ID)
	}

	jsonBytes, err := json.Marshal(team)
	if err != nil {
		http.Error(w, "Error encoding team to JSON: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error encoding team to JSON:", err)
		return
	}
	// fmt.Println(jsonBytes)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func GetMemberByIDHanlder(w http.ResponseWriter, r *http.Request) {
	memberIDBytes := r.URL.Path[len("/member/"):]
	memberID := string(memberIDBytes)
	fmt.Println(memberID)

	if len(memberID) == 0 {
		http.Error(w, "Member ID not provided", http.StatusBadRequest)
		return
	}
	objID, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		http.Error(w, "Invalid Member ID", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objID}

	var member model.HackMember
	err = db.HackMemberCollection.FindOne(context.Background(), filter).Decode(&member)
	if err != nil {
		http.Error(w, "Error retrieving member: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error retrieving member:", err)
		return
	}

	jsonBytes, err := json.Marshal(member)
	if err != nil {
		http.Error(w, "Error encoding member to JSON: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error encoding member to JSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func waitForTerminationSignal(server *http.Server) {
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
