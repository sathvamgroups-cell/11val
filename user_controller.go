package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"11val/models"
)

var validate = validator.New()

func validationErrorResponse(c *gin.Context, err error) {

	validationErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	errors := make(map[string]string)

	for _, validationError := range validationErrors {

		field := strings.ToLower(validationError.Field())
		tag := validationError.Tag()

		if tag == "required" {
			errors[field] = field + " is required"
			continue
		}

		if tag == "email" {
			errors[field] = field + " must be a valid email address"
			continue
		}

		if tag == "min" {
			errors[field] = field + " must be at least " + validationError.Param() + " characters"
			continue
		}

		if tag == "max" {
			errors[field] = field + " must be at most " + validationError.Param() + " characters"
			continue
		}

		if tag == "gte" {
			errors[field] = field + " must be at least " + validationError.Param()
			continue
		}

		if tag == "lte" {
			errors[field] = field + " must be at most " + validationError.Param()
			continue
		}

		errors[field] = "Invalid value for " + field
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"errors": errors,
	})
}

func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	if err := validate.Struct(user); err != nil {
		validationErrorResponse(c, err)
		return
	}

	models.DB.Create(&user)

	c.JSON(http.StatusCreated, user)
}

func GetUsers(c *gin.Context) {
	var users []models.User

	models.DB.Find(&users)

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	if result := models.DB.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	if result := models.DB.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	if err := validate.Struct(user); err != nil {
		validationErrorResponse(c, err)
		return
	}

	models.DB.Save(&user)

	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	if result := models.DB.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	models.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
