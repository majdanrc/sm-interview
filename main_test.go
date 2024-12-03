package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qwark97/interview/controller"
	"github.com/qwark97/interview/fetcher"
	"github.com/qwark97/interview/fetcher/model"
	"github.com/qwark97/interview/store"
)

func TestMain(m *testing.M) {
	// given
	var i int
	var link *string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if i == 0 {
			json.NewEncoder(writer).Encode(model.Response{
				NextLink: *link,
				Users: []model.User{
					{ID: "1", FirstName: "John", LastName: "Doe"},
					{ID: "2", FirstName: "Jane", LastName: "Doe"},
				},
			})
		} else if i == 1 {
			writer.WriteHeader(http.StatusTooManyRequests)
		} else if i == 2 {
			json.NewEncoder(writer).Encode(model.Response{
				NextLink: *link,
				Users: []model.User{
					{ID: "3", FirstName: "Jan", LastName: "Kowalski"},
					{ID: "4", FirstName: "Józef", LastName: "Nowak"},
				},
			})
		} else {
			json.NewEncoder(writer).Encode(model.Response{
				NextLink: "",
				Users: []model.User{
					{ID: "5", FirstName: "Marcin", LastName: "Plata"},
				},
			})
		}
		i++
	}))
	defer server.Close()
	link = &server.URL

	storage := store.New()
	fetcher := fetcher.New(server.URL, http.DefaultClient)

	controller := controller.New(storage, fetcher)

	// when
	controller.Handle(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	// then
}
