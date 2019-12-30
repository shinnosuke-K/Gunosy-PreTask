package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/shinnosuke-K/Gunosy-PreTask/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Information struct {
	DB *gorm.DB
}

type ResponseUserInfo map[string]string

func (info *Information) CreateHandler(c *gin.Context) {
	var account model.AccountInfo
	c.BindJSON(&account)

	if c.Request.ContentLength == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Body is missing",
			"cause":   "Body is empty",
		})
		return
	}

	inserted, msg := account.Insert(info.DB)
	if msg != "" {
		var cause string
		switch msg {
		case "required":
			cause = "required user_id and password"
		case "length":
			cause = "check the length of user_id and password"
		case "pattern":
			cause = ""
		case "duplication":
			cause = "already same user_id is used"
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Account creation failed",
			"cause":   cause,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account successfully created",
		"user": ResponseUserInfo{
			"user_id":  inserted.UserId,
			"nickname": inserted.Nickname,
		},
	})

}

func (info *Information) GetInfoHandler(c *gin.Context) {

	header := c.Request.Header.Get("Authorization")
	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Authentication Faild",
		})
		return
	}

	byteHeader, err := base64.StdEncoding.DecodeString(strings.Split(header, "Basic ")[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Authentication Faild",
		})
		return
	}

	idPass := strings.Split(string(byteHeader), ":")
	var account model.AccountInfo
	account.UserId = idPass[0]
	account.Password = idPass[1]

	getAccount, msg := account.AccountByID(info.DB)
	if msg != "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": msg,
		})
		return
	}

	responseBody := ResponseUserInfo{}
	rv := reflect.ValueOf(*getAccount)
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := fmt.Sprintf("%v", rv.FieldByName(field.Name))
		if tag := field.Tag.Get("json"); value != "" && tag != "password" {
			responseBody[field.Tag.Get("json")] = value
		}
	}

	if _, ok := responseBody["nickname"]; !ok {
		responseBody["nickname"] = responseBody["user_id"]
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User details by user_id",
		"user":    responseBody,
	})
	return
}

func (info *Information) UpdateHandler(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Authentication Faild",
		})
		return
	}

	byteHeader, err := base64.StdEncoding.DecodeString(strings.Split(header, "Basic ")[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Authentication Faild",
		})
		return
	}

	var account model.AccountInfo
	c.BindJSON(&account)
	if account.UserId != "" || account.Password != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User updating failed",
			"cause":   "not updatable user_id and password",
		})
		return
	}

	if account.Nickname == "" && account.Comment == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User updating failed",
			"cause":   "required nickname or comment",
		})
		return
	}

	idPass := strings.Split(string(byteHeader), ":")
	if idPass[0] != c.Param("user_id") {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "No Permission for Update",
		})
		return
	}

	account.UserId = idPass[0]
	account.Password = idPass[1]

	updated, msg := account.Update(info.DB)
	if msg != "" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully updated",
		"recipe": ResponseUserInfo{
			"nickname": updated.Nickname,
			"comment":  updated.Comment,
		},
	})
	return
}

func (info *Information) DeleteHandler(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Authentication Faild",
		})
		return
	}

	byteHeader, err := base64.StdEncoding.DecodeString(strings.Split(header, "Basic ")[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Authentication Faild",
		})
		return
	}

	idPass := strings.Split(string(byteHeader), ":")
	var account model.AccountInfo
	account.UserId = idPass[0]
	account.Password = idPass[1]

	deleted := account.Delete(info.DB)
	c.JSON(http.StatusOK, gin.H{
		"message": deleted,
	})
}
