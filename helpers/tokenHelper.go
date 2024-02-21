package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Prasenjit43/golang-jwt-project/database"
	// jwt "github.com/go-kit/kit/auth/jwt"
	// "github.com/golang-jwt/jwt/v5"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	User_type  string
	User_id    string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// var SECRET_KEY = os.Getenv("SECRET_KEY")
var SECRET_KEY string

func GenerateAllTokens(email string, firstName string, lastName string, userType string, userId string) (signedToken string, signedRefreshToken string, err error) {
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		log.Fatal(err)
	}
	SECRET_KEY = os.Getenv("SECRET_KEY")

	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		User_type:  userType,
		User_id:    userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, nil
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateObj primitive.D

	fmt.Println("User Id in updatealltoken :", userId)

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", updated_at})

	upsert := true
	// filter := bson.M{"user_id": userId}
	filter := bson.D{{"user_id", userId}}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	//update := bson.M{"$set": updateObj}

	updateResult, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)
	fmt.Println("Update Result :", updateResult)

	if err != nil {
		fmt.Println("Error in updateAllToken :", err.Error())
		log.Panic(err)
		return err
	}
	return nil

}

func UpdateAllTokens1(signedToken string, signedRefreshToken string, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// var updateObj primitive.D

	// updateObj = append(updateObj, bson.E{"token", signedToken})
	// updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	// updateObj = append(updateObj, bson.E{"updated_at", updatedAt})

	// upsert := true
	// filter := bson.M{"user_id": userId}
	// opt := options.UpdateOptions{
	// 	Upsert: &upsert,
	// }

	// update := bson.M{"$set": updateObj}

	fmt.Println("User Id in helper :", userId)

	filter := bson.D{{"user_id", userId}}
	update := bson.D{{"$set",
		bson.D{
			{"token", signedToken},
			{"refresh_token", signedRefreshToken},
			{"updated_at", updatedAt}},
	}}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		// bson.D{
		// 	{"$set": updateObj},
		// },
		update,
		// &opt,
	)
	fmt.Println("errr in helper :", err.Error())

	if err != nil {
		fmt.Println("errr1111 in helper :", err.Error())
		log.Panic(err)
		return err
	}
	return nil

}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		msg = fmt.Sprintf("invlaid token")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("invlaid token")
		msg = err.Error()
		return
	}

	return claims, msg

}
