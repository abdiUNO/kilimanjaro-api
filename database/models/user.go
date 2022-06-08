package models

import (
	"database/sql/driver"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"kilimanjaro-api/config"
	"kilimanjaro-api/database/orm"
	"kilimanjaro-api/utils"
	"strings"
)

/*
JWT claims struct
*/

var cfg = config.GetConfig()

type Token struct {
	UserId string
	jwt.StandardClaims
}

type UserType string

const (
	BUYER  UserType = "buyer"
	SELLER UserType = "seller"
)

func (ut *UserType) Scan(value interface{}) error {
	*ut = UserType(value.([]byte))
	return nil
}

func (ut UserType) Value() (driver.Value, error) {
	if len(ut) == 0 {
		return nil, nil
	}
	return string(ut), nil
}

type User struct {
	orm.GormModel
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Email         string    `json:"email"`
	Image         string    `json:"image"`
	Phone         string    `json:"phone"`
	Location      string    `json:"location"`
	Website       string    `json:"website"`
	JwtToken      string    `sql:"-" json:"jwtToken"`
	UserType      UserType  `json:"userType" sql:"type:ENUM('buyer','seller');DEFAULT:null"`
	EmailVerified bool      `sql:"not null;DEFAULT:false" json:"emailVerified"`
	Secret        string    `json:"-"`
	Products      []Product `json:"products"`
}

func (user *User) TableName() string {
	return "users"
}

func (user *User) Validate() *utils.Error {
	if !strings.Contains(user.Email, "@") {
		return utils.NewError(utils.EINVALID, "email address is required", nil)
	}

	temp := &User{}

	err := GetDB().Table("users").Where("email = ?", user.Email).First(temp).Error
	//fmt.Println(temp == nil)

	if err != nil && err != gorm.ErrRecordNotFound {
		return utils.NewError(utils.EINVALID, "DB error: user record not found", err)
	}

	if temp.Email != "" && temp.Email == user.Email {
		return utils.NewError(utils.EINVALID, "email address already in use by another user", nil)
	}

	return nil
}

func (user *User) Create() (*User, *utils.Error) {
	err := GetDB().Create(user).Error
	if err != nil {
		log.Error(err)
		return &User{}, utils.NewError(utils.ECONFLICT, "DB error: could not create user", nil)
	}

	//Create new JWT token for the newly registered account
	tk := &Token{
		UserId: user.ID,
		//StandardClaims: jwt2.StandardClaims{ExpiresAt: 150000},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))
	user.JwtToken = tokenString //Store the token in the response

	return user, nil
}

func ValidateUserInfo(id string, user *User) *utils.Error {
	if !strings.Contains(user.Email, "@") {
		return utils.NewError(utils.EINVALID, "email address is required", nil)

	}

	temp := &User{}
	err := GetDB().Table("users").Where("id = ?", id).First(temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return utils.NewError(utils.EINVALID, "DB error: user record not found", err)
	}

	//log.Println(temp.Email != user.Email)

	if user.Email != "" && temp.Email != user.Email {
		tempUserEmail := &User{}
		tempErr := GetDB().Table("users").Where("email = ?", user.Email).First(tempUserEmail).Error
		fmt.Println(tempUserEmail.Email)
		if tempErr != gorm.ErrRecordNotFound && tempUserEmail.Email == user.Email {
			return utils.NewError(utils.EINVALID, "DB error: email address already in use by another user", nil)
		}
	}

	return nil
}

func Update(id string, user *User) (*User, *utils.Error) {

	if validateErr := ValidateUserInfo(id, user); validateErr != nil {
		return &User{}, validateErr
	}

	updateUser, err := FindUserById(id)

	//Create JWT token
	tk := &Token{UserId: updateUser.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))
	updateUser.JwtToken = tokenString //Store the token in the response

	err = GetDB().Model(&updateUser).Updates(&user).Error
	if err != nil {
		return &User{}, utils.NewError(utils.ECONFLICT, "DB error: could not update user", nil)
	}

	return updateUser, nil
}

func Login(email string) (*User, *utils.Error) {

	user := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &User{}, utils.NewError(utils.ENOTFOUND, "DB error: user record not found", nil)
		}

		return &User{}, utils.NewError(utils.EINTERNAL, "internal server error", nil)
	}

	//Create JWT token
	tk := &Token{
		UserId: user.ID,
		//StandardClaims: jwt2.StandardClaims{ExpiresAt: 150000},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))
	user.JwtToken = tokenString //Store the token in the response

	return user, nil
}

func QueryUsers(userID string, query string) (*[]User, *utils.Error) {
	users := &[]User{}
	var blockedList []string
	var idStr string

	err := GetDB().Table("blockeds").Where("user_id = ?", userID).Pluck("friend_id", &blockedList).Error

	if err != nil {
		return &[]User{}, utils.NewError(utils.EINVALID, "invalid login credentials. Please try again", err)
	}

	for i, id := range blockedList {
		if i == 0 {
			idStr += "'" + id + "'"
		} else {
			idStr += ",'" + id + "'"
		}
	}

	fmt.Println(len(idStr))

	if len(idStr) <= 0 {
		err = GetDB().Table("users").Where("name LIKE ?", query+"%").Find(&users).Error
	} else {
		err = GetDB().Table("users").Where("id NOT IN ("+idStr+") AND name LIKE ?", query+"%").Find(&users).Error
	}

	if err != nil {
		return &[]User{}, utils.NewError(utils.EINVALID, "invalid login credentials. Please try again", err)
	}

	return users, nil
}

func FindUserById(u string) (*User, error) {

	user := &User{}
	err := GetDB().Table("users").Where("id = ?", u).First(user).Error

	if user.Email == "" { //User not found!
		return nil, err
	}

	return user, err
}

func GenerateJwtToken(user *User) string {

	tk := &Token{
		UserId: user.ID,
		//StandardClaims: jwt2.StandardClaims{ExpiresAt: 150000},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))

	return tokenString
}

func FindUserByEmail(email string) (*User, error) {

	user := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(user).Error

	if user.Email == "" { //User not found!
		return nil, err
	}

	return user, err
}
