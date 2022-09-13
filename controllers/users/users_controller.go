package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/deeptithakur22/bookstore_oauth-go/oauth"
	"github.com/deeptithakur22/bookstore_users-api/domain/users"
	"github.com/deeptithakur22/bookstore_users-api/services"
	"github.com/deeptithakur22/bookstore_utils-go/rest_errors"

	"github.com/gin-gonic/gin"
)

func getUserId(userIdParam string) (int64, *rest_errors.RestErr) {
	userID, userErr := strconv.ParseInt(userIdParam, 10, 64)
	if userErr != nil {
		return 0, rest_errors.NewBadRequestError("user id should be a number")

	}
	return userID, nil
}

func Create(c *gin.Context) {
	var user users.User
	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		//TODO: Handle error
		fmt.Println("Error while reading request body", err.Error())
		return
	}
	if err := json.Unmarshal(bytes, &user); err != nil {
		//TODO: Handle json error
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		fmt.Println("error on unmarshaling: ", err.Error())
		return
	}

	fmt.Println(user)
	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		//TODO: Handle User Creation Error
		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Get(c *gin.Context) {
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status, err)
		return
	}
	userID, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}
	user, getErr := services.UsersService.GetUser(userID)
	if getErr != nil {
		//TODO: Handle User Creation Error
		c.JSON(getErr.Status, getErr)
		return
	}

	if oauth.GetCallerId(c.Request) == user.Id {
		c.JSON(http.StatusOK, user.Marshall(false))
	}
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))

}

func Update(c *gin.Context) {
	userID, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	var user users.User
	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		//TODO: Handle error
		fmt.Println("Error while reading request body", err.Error())
		return
	}
	if err := json.Unmarshal(bytes, &user); err != nil {
		//TODO: Handle json error
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		fmt.Println("error on unmarshaling: ", err.Error())
		return
	}
	user.Id = userID
	isPartial := c.Request.Method == http.MethodPatch
	result, Err := services.UsersService.UpdateUser(isPartial, user)
	if err != nil {
		c.JSON(Err.Status, Err)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Delete(c *gin.Context) {
	userID, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}
	if err := services.UsersService.DeleteUser(userID); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "Deleted"})
}

func Search(c *gin.Context) {
	status := c.Query("status")
	users, err := services.UsersService.SearchUser(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, users.Marshall(c.GetHeader("X-Public") == "true"))
}

func Login(c *gin.Context) {
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	user, err := services.UsersService.LoginUser(request)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
}
