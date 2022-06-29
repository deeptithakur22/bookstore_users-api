package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/deeptithakur22/bookstore_users-api/domain/users"
	"github.com/deeptithakur22/bookstore_users-api/services"
	"github.com/deeptithakur22/bookstore_users-api/utils/errors"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user users.User
	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		//TODO: Handle error
		fmt.Println("Error while reading request body", err.Error())
		return
	}
	if err := json.Unmarshal(bytes, &user); err != nil {
		//TODO: Handle json error
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		fmt.Println("error on unmarshaling: ", err.Error())
		return
	}

	fmt.Println(user)
	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		//TODO: Handle User Creation Error
		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func GetUser(c *gin.Context) {
	userID, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("user id should be a number")
		c.JSON(err.Status, err)
		return
	}
	user, getErr := services.GetUser(userID)
	if getErr != nil {
		//TODO: Handle User Creation Error
		c.JSON(getErr.Status, getErr)
		return
	}

	c.JSON(http.StatusOK, user)

}

func SearchUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "implement me!")

}
