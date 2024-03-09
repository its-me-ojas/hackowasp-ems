package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"hackowasp_ems/db"
	"hackowasp_ems/model"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllMembersHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve all members from the database
	members, err := db.HackMemberCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Error retrieving members: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error retrieving members:", err)
		return
	}

	var allMembers []model.HackMember
	for members.Next(context.Background()) {
		var member model.HackMember
		err := members.Decode(&member)
		if err != nil {
			http.Error(w, "Error decoding member: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error decoding member:", err)
			return
		}
		allMembers = append(allMembers, member)
	}

	jsonBytes, err := json.Marshal(allMembers)
	if err != nil {
		http.Error(w, "Error encoding members to JSON: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error encoding members to JSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

func GetAllTeamsHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve all teams from the database
	teams, err := db.HackTeamCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Error retrieving teams: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error retrieving teams:", err)
		return
	}

	var allTeams []model.HackTeam
	for teams.Next(context.Background()) {
		var team model.HackTeam
		err := teams.Decode(&team)
		if err != nil {
			http.Error(w, "Error decoding team: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error decoding team:", err)
			return
		}
		allTeams = append(allTeams, team)
	}

	jsonBytes, err := json.Marshal(allTeams)
	if err != nil {
		http.Error(w, "Error encoding teams to JSON: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error encoding teams to JSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
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
	currentTime := time.Now()
	fmt.Println("Current Time:", currentTime)

	member.Time = currentTime

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

	teamID := member.TeamID
	if member.CheckedIn {
		teamUpdate := bson.M{
			"$inc": bson.M{"checkedInCount": 1},
			// "$set": bson.M{"checkedIn": true},
		}
		// Update the team's checkedInCount
		_, err = db.HackTeamCollection.UpdateOne(context.Background(), bson.M{"_id": teamID}, teamUpdate)
	} else {
		teamUpdate := bson.M{
			"$inc": bson.M{"checkedInCount": -1},
			// "$set": bson.M{"checkedIn": false},
		}
		// Update the team's checkedInCount
		_, err = db.HackTeamCollection.UpdateOne(context.Background(), bson.M{"_id": teamID}, teamUpdate)

	}
	if err != nil {
		http.Error(w, "Error updating team: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error updating team:", err)
		return
	}

	fmt.Println("Member check-in status : ", member.CheckedIn)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Member check-in ", "status": "` + fmt.Sprint(member.CheckedIn) + `"}`))
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
