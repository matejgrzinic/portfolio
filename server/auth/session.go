package auth

import (
	"time"

	"github.com/matejgrzinic/portfolio/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	Agent   string             `json:"agent"`
	Address string             `json:"address"`
	SID     string             `json:"sid"`
	UserID  primitive.ObjectID `json:"userid" bson:"_id,omitempty"`
	Expires int64              `json:"expires"`
}

func SessionBySID(SID string, address string, agent string, db db.API) (*Session, error) {
	var s Session
	err := db.QueryRow(
		"session by id",
		"session",
		bson.M{
			"agent":   agent,
			"address": address,
			"sid":     SID,
			"expires": bson.M{"$gt": time.Now().Unix()},
		},
		nil,
		&s,
	)
	return &s, err
}
