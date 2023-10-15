package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID		int		`json:"id"`
	Name	string	`json:"name"`
	Email 	string	`json:"email"`
	Age		int		`json:"age"`
}

var users = []user{
	{ ID: 1, Name: "Rohit", Email: "rohit@mail.com", Age: 21 },
	{ ID: 2, Name: "Vaibhav", Email: "vaibhav@mail.com", Age: 23 },
}

func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

func createUser(c *gin.Context) {
	var newUser user

	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{ "message": "Something went wrong!" })
		return
	}

	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func getUserbyId(id int) (*user, error) {
	for idx, user := range users {
		if (user.ID == id) {
			return &users[idx], nil
		}
	}

	return nil, errors.New("User not found!")
}

func userById(c *gin.Context) {
	id := c.Param("id")
	idNum, err := strconv.Atoi(id);

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{ "message": "Invalid user ID!" })
		return
	}

	user, err := getUserbyId(idNum)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{ "message": "User not found!" })
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func updateUser(c *gin.Context) {
	var userToUpdate user

	if err := c.BindJSON(&userToUpdate); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{ "message": "Something went wrong!" })
		return
	}

	user, err := getUserbyId(userToUpdate.ID)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{ "message": "User not found!" })
		return
	}

	user.Name = userToUpdate.Name
	user.Age = userToUpdate.Age
	user.Email = userToUpdate.Email
	
	c.IndentedJSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	idNum, err := strconv.Atoi(id);

	user, err := getUserbyId(idNum)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{ "message": "User not found!" })
		return
	}

	users = filterUser(users, user.ID)
	c.IndentedJSON(http.StatusOK, user)
}

func filterUser(slice []user, id int) []user {
	var res []user

	for idx, user := range slice {
		if (user.ID != id) {
			res = append(res, slice[idx])
		}
	}

	return res
}

func main() {
	router := gin.Default()
	router.GET("/users", getUsers)
	router.GET("/users/:id", userById)
	router.POST("/users", createUser)
	router.PATCH("/users", updateUser)
	router.DELETE("/users/:id", deleteUser)

	router.Run("localhost:8080")
}