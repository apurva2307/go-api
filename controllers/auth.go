package controllers

import (
	"time"

	"github.com/apurva2307/go-api/db"
	"github.com/apurva2307/go-api/helpers"
	"github.com/apurva2307/go-api/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type SingnUpUserBody struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=7,max=50"`
}
type LoginUserBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=7,max=50"`
}

func SignupUser(c *fiber.Ctx) error {
	collection := db.Mg.Db.Collection("users")
	userBody := new(SingnUpUserBody)
	err := c.BodyParser(userBody)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": []map[string]string{{"msg": "Kindly provide proper inputs"}},
		})
	}
	errors := helpers.ValidateStruct(userBody)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}
	var userType *models.User
	err = collection.FindOne(c.Context(), bson.D{{Key: "email", Value: userBody.Email}}).Decode(&userType)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": map[string]string{"msg": "User with provided email already exists."},
		})
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(userBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": err.Error()},
		})
	}
	user := &models.User{
		Name:      userBody.Name,
		Email:     userBody.Email,
		Password:  hashed,
		CreatedAt: time.Now(),
	}
	result, err := collection.InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": err.Error()},
		})
	}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}
	var createdUser models.User
	err = collection.FindOne(c.Context(), query).Decode(&createdUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": err.Error()},
		})
	}
	token, err := helpers.CreateJwtToken(createdUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": err.Error()},
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"user":  createdUser,
		"token": token,
	})
}

func LoginUser(c *fiber.Ctx) error {
	collection := db.Mg.Db.Collection("users")
	userBody := new(LoginUserBody)
	err := c.BodyParser(userBody)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": []map[string]string{{"msg": "Kindly provide proper inputs"}},
		})
	}
	errors := helpers.ValidateStruct(userBody)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}
	var user *models.User
	err = collection.FindOne(c.Context(), bson.D{{Key: "email", Value: userBody.Email}}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": map[string]string{"msg": "Invalid credentials provided."},
		})
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(userBody.Password))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": "Invalid credentials provided."},
		})
	}
	token, err := helpers.CreateJwtToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": err.Error()},
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": "user successfully logged in.",
		"token":   token,
	})
}
