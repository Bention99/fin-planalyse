package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (a *app) getCurrentUserID(r *http.Request) (uuid.UUID, error) {
	c, err := r.Cookie("session")
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(c.Value)
}