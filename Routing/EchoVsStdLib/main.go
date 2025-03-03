package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/user", handleEchoGetUser)


}

type APIError struct {
	Status int
	Message string
}

func handleEchoGetUser(c echo.Context) error {
	user, err := getUser()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)

}

type User struct {}

func getUser() (*User, error) {
	return nil, nil
}
