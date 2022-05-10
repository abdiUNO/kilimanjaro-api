package models

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"kilimanjaro-api/config"
	"kilimanjaro-api/database/orm"
	"kilimanjaro-api/utils"
	"log"
	"strings"
)

/*
JWT claims struct
*/

type Token struct {
	UserId string
	jwt.StandardClaims
}

var cfg = config.GetConfig()

type User struct {
	orm.GormModel
	Name          string `json:"name"`
	Description   string `json:"description"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Password      string `json:"-"`
	Location      string `json:"location"`
	Website       string `json:"website"`
	JwtToken      string `sql:"-" json:"jwtToken"`
	EmailVerified bool   `sql:"not null;DEFAULT:false" json:"emailVerified"`
	Secret        string `json:"-"`
	Vendor        Vendor `json:"-"`
}

func (user *User) TableName() string {
	return "users"
}

func (user *User) Validate() *utils.Error {
	if !strings.Contains(user.Email, "@") {
		return utils.NewError(utils.EINVALID, "email address is required", nil)
	}

	//if len(user.Password) < 6 {
	//	return utils.NewError(utils.EINVALID, "password is required", nil)
	//}

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
	//hashedPassword := hashAndSalt([]byte(user.Password))
	//user.Password = string(hashedPassword)

	user.Password = ""

	err := GetDB().Create(user).Error
	if err != nil {
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

	user.Password = "" //remove password

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

	log.Println(temp.Email != user.Email)

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
		return &User{}, utils.NewError(utils.ECONFLICT, "DB error: could not create user", nil)
	}

	return updateUser, nil
}

func (user *User) UpdatePassword(oldPassword string, newPassword string) *utils.Error {

	newPasswordHashed := hashAndSalt([]byte(newPassword))

	fmt.Println(newPassword)
	fmt.Println(user.Password)

	if comparePasswords(user.Password, []byte(oldPassword)) == false { //Password does not match!
		return utils.NewError(utils.EINVALID, "Incorrect current password", nil)
	}

	if len(newPassword) < 6 {
		return utils.NewError(utils.EINVALID, "Password is required", nil)
	}

	user.Password = string(newPasswordHashed)

	//Create JWT token
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))
	user.JwtToken = tokenString //Store the token in the response

	err := GetDB().Save(&user).Updates(&user).Error
	if err != nil {
		return utils.NewError(utils.ECONFLICT, "DB error: could not update password", nil)
	}

	return nil
}

func Login(email string) (*User, *utils.Error) {

	user := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &User{}, utils.NewError(utils.ENOTFOUND, "Your email or password is incorrect.", nil)
		}

		return &User{}, utils.NewError(utils.EINTERNAL, "internal server error", nil)
	}

	//if comparePasswords(user.Password, []byte(password)) == false { //Password does not match!
	//	return &User{}, utils.NewError(utils.EINVALID, "Your email or password is incorrect.", nil)
	//}

	//Worked! Logged In
	user.Password = ""

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
		err = GetDB().Table("users").Where("full_name LIKE ?", query+"%").Find(&users).Error
	} else {
		err = GetDB().Table("users").Where("id NOT IN ("+idStr+") AND full_name LIKE ?", query+"%").Find(&users).Error
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

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
