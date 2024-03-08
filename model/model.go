package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HackMember struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	TeamID      primitive.ObjectID `json:"teamID"`
	MemberName  string             `json:"memberName"`
	MemberEmail string             `json:"memberEmail"`
	MemberPhone int                `json:"memberPhone"`
	MemberRoll  int                `json:"memberRoll"`
	CheckedIn   bool               `json:"checkedIn"`
	Time        time.Time          `json:"time"`
	V           int                `json:"__v"`
}

type HackTeam struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	TeamName       string             `json:"teamName"`
	College        string             `json:"college"`
	LeaderEmail    string             `json:"leaderEmail"`
	Members        []HackMember       `json:"members"`
	CheckedIn      bool               `json:"checkedIn"`
	CheckedInCount int                `json:"checkedInCount"`
	TotalMembers   int                `json:"totalMembers"`
	LeaderName     string             `json:"leaderName"`
	Selected       bool               `json:"selected"`
	Locked         bool               `json:"locked"`
	Waitlist       bool               `json:"waitlist"`
}
