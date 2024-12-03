package store

import (
	"errors"
	"fmt"

	"github.com/qwark97/interview/store/model"
)

var ErrDuplicate = errors.New("entity already exists")

type Store struct {
}

func New() Store {
	return Store{}
}

func (s Store) InsertUser(user model.User) error {
	/*
		...
	*/
	fmt.Printf("Inserting user: %v\n", user)

	return nil
}
