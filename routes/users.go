package routes

import (
	"fmt"
	"net/http"
	"voiceassistant/models"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func UsersLogin(c *gin.Context) {
	user := models.User{}
	//err := c.ShouldBindJSON(&user)
	user.Email = c.Request.PostFormValue("email")
	user.Password = c.Request.PostFormValue("password")

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	err := user.IsAuthenticated(&conn)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := user.GetAuthToken()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"status":"success",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error": "There was an error authenticating.",
		"status": "failed",
	})
}

func UsersRegister(c *gin.Context) {
	user := models.User{}
	//err := c.ShouldBindJSON(&user)

	user.FullName = c.Request.PostFormValue("fullname")
	user.Email = c.Request.PostFormValue("email")
	user.Password = c.Request.PostFormValue("password")
	user.PasswordConfirm = c.Request.PostFormValue("password_confirm")

	PhoneNo, _ := strconv.Atoi(c.Request.PostFormValue("phoneno"))
	user.PhoneNo = uint64(PhoneNo)
	user.StoreName = c.Request.PostFormValue("storename")
	user.StoreAddress = c.Request.PostFormValue("storeaddress")
	PinCode, _ := strconv.Atoi(c.Request.PostFormValue("pincode"))
	user.PinCode = int64(PinCode)
    fmt.Println(user)
	
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	err := user.Register(&conn)
	if err != nil {
		fmt.Println("Error in user.Register()")
		c.JSON(http.StatusBadRequest , gin.H{"error": err.Error()})
		return
	}

	token, err := user.GetAuthToken()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"status":"success",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token" : token,
		"user_id": user.ID,
		"status" : "success",

	})
}

