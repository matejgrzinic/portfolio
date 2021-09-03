package portfolio

import (
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name string             `json:"name"`
}

func (p *Portfolio) UserByName(name string) (*User, error) {
	user := new(User)
	err := p.DB().QueryRow(
		"user by name",
		"users",
		bson.D{{Key: "user", Value: user}},
		nil,
		user,
	)

	return user, err
}

func (p *Portfolio) UserByID(ID primitive.ObjectID) (*User, error) {
	user := new(User)
	err := p.DB().QueryRow(
		"user by id",
		"users",
		bson.D{{Key: "_id", Value: ID}},
		nil,
		user,
	)

	return user, err
}

func (p *Portfolio) AllUsers() ([]User, error) {
	result := make([]User, 0)
	user := new(User)
	err := p.DB().QueryRows(
		"user by name",
		"users",
		bson.D{},
		nil,
		user,
		func() error {
			var cpyU User
			copier.CopyWithOption(&cpyU, user, copier.Option{DeepCopy: true})
			result = append(result, cpyU)
			return nil
		},
	)

	return result, err
}
