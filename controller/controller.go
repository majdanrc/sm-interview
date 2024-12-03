package controller

import (
	"context"
	"fmt"
	"github.com/qwark97/interview/fetcher"
	"net/http"
	"time"

	fetcherModel "github.com/qwark97/interview/fetcher/model"
	storeModel "github.com/qwark97/interview/store/model"
)

type Store interface {
	InsertUser(user storeModel.User) error
}

type Controller struct {
	storage Store
	fetcher fetcher.Fetcher
}

func New(storage Store, fetcher fetcher.Fetcher) Controller {
	return Controller{storage: storage, fetcher: fetcher}
}

func (c Controller) Handle(writer http.ResponseWriter, request *http.Request) {
	/* TODO:
	- Fetch users from vendor API
	- Insert new users into database
	- In case of processing failure (unrecoverable), return 500
	*/

	ctx, cancel := context.WithTimeout(request.Context(), time.Minute)
	defer cancel()

	processor := func(user fetcherModel.User) error {
		storeUser := storeModel.User{
			ID:       user.ID,
			FullName: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		}

		return c.storage.InsertUser(storeUser)
	}

	err := c.fetcher.Users(ctx, processor)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}
