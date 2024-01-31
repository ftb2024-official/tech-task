package handlers

import (
	"encoding/json"
	"os"
	"time"

	util "tech_task/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
)

type Country struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type User struct {
	tableName     struct{}  `pg:"users, alias:u"`
	ID            string    `pg:"id, pk, type:uuid, notnull" json:"id"`
	Name          string    `pg:"name, type:text, notnull" json:"name"`
	Surname       string    `pg:"surname, type:text, notnull" json:"surname"`
	Patronymic    string    `pg:"patronymic, type:text" json:"patronymic"`
	Age           int       `pg:"age, type:integer, notnull" json:"age"`
	Gender        string    `pg:"gender, type:varchar(7), notnull" json:"gender"`
	Countries     []Country `pg:"countries, type:jsonb[], notnull" json:"country"`
	UserCreatedAt time.Time `pg:"user_created_at, type:timestamptz, notnull" json:"user_created_at"`
}

func GetUsersHandler(db *pg.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var users []User

		query := db.Model(&users)
		limit, _ := ctx.Get("limit")
		query.Limit(limit.(int))

		if err := query.Select(); err != nil {
			ctx.JSON(500, gin.H{"error": "Internal Server Error."})
			return
		}

		ctx.JSON(200, gin.H{"data": users})
	}
}

func CreateUserHandler(db *pg.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newUser User

		name, _ := ctx.Get("name")
		surname, _ := ctx.Get("surname")
		newUser.Surname = surname.(string)

		ageURL := os.Getenv("AGIFY_API_URL") + name.(string)
		age, err := util.FetchData(ageURL)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to fetch data from external API"})
			return
		}

		genderURL := os.Getenv("GENDERIZE_API_URL") + name.(string)
		gender, err := util.FetchData(genderURL)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to fetch data from external API"})
			return
		}

		nationalityURL := os.Getenv("NATIONALIZE_API_URL") + name.(string)
		nation, err := util.FetchData(nationalityURL)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to fetch data from external API"})
			return
		}

		if err := json.Unmarshal(age, &newUser); err != nil {
			logrus.Printf("Unable to marshal age due to (%v)", err)
		}

		if err := json.Unmarshal(gender, &newUser); err != nil {
			logrus.Printf("Unable to marshal gender due to (%v)", err)
		}

		type NationalityResponse struct {
			Count   int       `json:"-"`
			Name    string    `json:"-"`
			Country []Country `json:"country"`
		}
		var natResponse NationalityResponse
		if err := json.Unmarshal(nation, &natResponse); err != nil {
			logrus.Printf("Unable to marshal countries (%v)", err)
		}
		newUser.Countries = natResponse.Country

		if _, err := db.Model(&newUser).Insert(); err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to create a new user"})
			return
		}

		ctx.JSON(201, gin.H{"data": newUser})
	}
}

func EditUserHandler(db *pg.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var oldUser, editedUser User

		userID, _ := ctx.Params.Get("id")

		if err := db.Model(&oldUser).Where("id = ?", userID).Select(); err != nil {
			ctx.JSON(404, gin.H{"error": "User not found"})
			return
		}

		if err := ctx.BindJSON(&editedUser); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid JSON format."})
			return
		}

		if editedUser.Name != "" {
			oldUser.Name = editedUser.Name
		}

		if editedUser.Surname != "" {
			oldUser.Surname = editedUser.Surname
		}

		if editedUser.Patronymic != "" {
			oldUser.Patronymic = editedUser.Patronymic
		}

		if editedUser.Gender != "" {
			oldUser.Gender = editedUser.Gender
		}

		if editedUser.Age != 0 {
			oldUser.Age = editedUser.Age
		}

		if editedUser.Countries != nil {
			oldUser.Countries = editedUser.Countries
		}

		if _, err := db.Model(&oldUser).Where("id = ?", userID).Update(); err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to update user."})
			return
		}

		ctx.JSON(200, gin.H{"data": oldUser})
	}
}

func DeleteUserHandler(db *pg.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user User
		id, _ := ctx.Params.Get("id")

		if err := db.Model(&user).Where("id = ?", id).Select(); err != nil {
			ctx.JSON(404, gin.H{"error": "User not found"})
			return
		}

		if _, err := db.Model(&user).Where("id = ?", id).Delete(); err != nil {
			ctx.JSON(500, gin.H{"error": err})
			return
		}

		ctx.JSON(200, gin.H{"success": "User successfully deleted"})
	}
}
