package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Prasenjit43/golang-jwt-project/database"
	"github.com/Prasenjit43/golang-jwt-project/helpers"
	"github.com/Prasenjit43/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var validate = validator.New()

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(hashedPassword)
}

func VerifyPassword(foundPassword string, userTypedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userTypedPassword), []byte(foundPassword))
	check := true
	msg := ""

	if err != nil {
		return false, "mismatched password"
	}
	return check, msg
}

func Signup() gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		fmt.Println("2222222")

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		fmt.Println("!!!!!!!! User Name : ", user.Email)
		fmt.Println("!!!!!!!! User Name : ", user.Email)
		fmt.Println("!!!!!!!! User Name : ", user.ID)

		validateErr := validate.Struct(user)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validateErr.Error(),
			})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			fmt.Print("Errr while checking email :", err.Error())
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "found error while checking for email",
			})

		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "record with this email is already present",
			})
		}

		hashedPassword := HashPassword(*user.Password)
		user.Password = &hashedPassword

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {

			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occured while checking for phone no",
			})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "record with this phone no is already present",
			})
			return
		}

		user.Created_at, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		user.Updated_at, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken, _ := helpers.GenerateAllTokens(
			*user.Email,
			*user.First_name,
			*user.Last_name,
			*user.User_type,
			user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNum, insertErr := userCollection.InsertOne(ctx, &user)
		if insertErr != nil {
			msg := "User item not created"
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": msg,
			})
			return
		}
		c.JSON(http.StatusOK, resultInsertionNum)

	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "email & password are incorrect",
			})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "email id not found",
			})
		}

		token, refreshtoken, _ := helpers.GenerateAllTokens(
			*foundUser.Email,
			*foundUser.First_name,
			*foundUser.Last_name,
			*foundUser.User_type,
			foundUser.User_id)
		fmt.Println("Generated Token :", token)
		fmt.Println("foundUser.User_id :", foundUser.User_id)

		helpers.UpdateAllTokens(token, refreshtoken, foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, foundUser)

	}

}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}

		groupStage := bson.D{{
			"$group", bson.D{
				{"_id", bson.D{{"_id", "null"}}},
				{"total_count", bson.D{{"$sum", 1}}},
				{"data", bson.D{{"$push", "$$ROOT"}}},
			}}}

		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}}}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		var allUsers []bson.M
		if err := result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		fmt.Println("UserID GU -  ", userId)

		if err := helpers.MatchUserTypeWithUserId(c, userId); err != nil {
			c.JSON(
				http.StatusBadRequest, gin.H{
					"error": err.Error(),
				},
			)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var users models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&users)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}
