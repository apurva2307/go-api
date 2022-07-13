package controllers

import (
	"fmt"

	"github.com/apurva2307/go-api/db"
	"github.com/apurva2307/go-api/helpers"
	"github.com/apurva2307/go-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllEmployees(c *fiber.Ctx) error {
	cursor, err := db.Mg.Db.Collection("employees").Find(c.Context(), bson.D{{}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": err.Error()},
		})
	}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	fmt.Println(claims["email"].(string))
	var employees []models.Employee = make([]models.Employee, 0)
	err = cursor.All(c.Context(), &employees)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": err.Error()},
		})
	}
	return c.JSON(employees)
}
func CreateEmployee(c *fiber.Ctx) error {
	collection := db.Mg.Db.Collection("employees")
	employee := new(models.Employee)
	err := c.BodyParser(employee)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": err.Error()},
		})
	}
	errors := helpers.ValidateStruct(employee)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, err := primitive.ObjectIDFromHex(claims["user_id"].(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": "userid is not valid."},
		})
	}
	employee.User = userID
	result, err := collection.InsertOne(c.Context(), employee)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": err.Error()},
		})
	}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}
	var createdEmployee *bson.M
	err = collection.FindOne(c.Context(), query).Decode(&createdEmployee)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": err.Error()},
		})
	}
	return c.Status(200).JSON(createdEmployee)
}

func UpdateEmployee(c *fiber.Ctx) error {
	collection := db.Mg.Db.Collection("employees")
	employeeId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": "Provided Id is not valid."},
		})
	}
	employee := new(models.Employee)
	if err := c.BodyParser(employee); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": err.Error()},
		})
	}
	errors := helpers.ValidateStruct(employee)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}
	query := bson.D{{Key: "_id", Value: employeeId}}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "name", Value: employee.Name},
			},
		},
	}
	err = collection.FindOneAndUpdate(c.Context(), query, update).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"errors": map[string]string{"msg": "No data found with provided id."},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": err.Error()},
		})
	}

	return c.Status(fiber.StatusOK).JSON(employee)
}

func DeleteEmployee(c *fiber.Ctx) error {
	collection := db.Mg.Db.Collection("employees")
	employeeId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": []map[string]string{{"msg": "Provided Id is not valid."}},
		})
	}
	query := bson.D{{Key: "_id", Value: employeeId}}
	result, err := collection.DeleteOne(c.Context(), &query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": []map[string]string{{"msg": err.Error()}},
		})
	}
	if result.DeletedCount < 1 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errors": []map[string]string{{"msg": "No data found with provided id."}},
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": "Record deleted successfully.",
	})

}
