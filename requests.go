package main

type RequestWithSession struct {
	SessionID string `json:"sessionID" validate:"required"`
}

type RequestUser struct {
	RequestWithSession
	Username string `json:"username" validate:"required"`
}
