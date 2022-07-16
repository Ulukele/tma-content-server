package main

type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}

type Team struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	OwnerId uint   `json:"ownerId"`
}

type Board struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	TeamId uint   `json:"teamId"`
}
