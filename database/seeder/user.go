package seeder

import (
	"user-service/constants"
	"user-service/domain/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserSeeder(db gorm.DB) {
	password, _ := bcrypt.GenerateFromPassword([]byte("P@ssw0rd123"), bcrypt.DefaultCost)
	users := models.User{
		UUID:        uuid.New(),
		Name:        "Admin",
		Username:    "admin",
		Password:    string(password),
		Email:       "admin@gmail.com",
		PhoneNumber: "0812131",
		RoleId:      constants.Admin,
	}

	err := db.FirstOrCreate(&users, models.User{Username: users.Username}).Error
	if err != nil {
		logrus.Errorf("failed to seed user: %v", err)
		panic(err)
	}

	logrus.Infof("user %s successfuly seeded", users.Username)
}
