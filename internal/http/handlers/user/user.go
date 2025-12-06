package user

import (
	// "crypto/sha256"
	// "encoding/hex"
	"encoding/json"
	"os"

	// "strings"
	"time"

	// "go/types"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/roshanansy/file-upload/internal/database"
	"github.com/roshanansy/file-upload/internal/model"
	"github.com/roshanansy/file-upload/internal/types"
	"golang.org/x/crypto/bcrypt"

	// "github.com/roshanansy/file-upload/internal/model"
	"github.com/roshanansy/file-upload/internal/utils/response"
)

// Define your User model — or import it from your models package

//  func CreateUser(w http.ResponseWriter, r *http.Request) {
// 	var user types.User
// 	// data :=r.Body
// 	// Decode request body
// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		response.WriteJson(w, http.StatusBadRequest, map[string]string{
// 			"error": "Invalid request payload",
// 		})
// 		return
// 	}

// 	//hash the password before storing it in database

// 	// bytePassword :=sha256.Sum256([]byte(user.Password))

// 	// passwordHash:=hex.EncodeToString(bytePassword[:])

// 	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 		if err != nil {
// 			response.WriteJson(w, http.StatusBadRequest, map[string]string{
// 			"error": "Invalid request payload",
// 		})
// 		return
// 		}

// 	modelUser := model.User{
// 		ID:       uuid.New().String(),
// 		Name:     user.Name,
// 		Email:    user.Email,
// 		Password: string(passwordHash), // In real applications, hash the password before storing
// 	}

// 	// return res:=database.DB.Create(&modelUser);

// 	// Save to database
// 	if err := database.DB.Create(&modelUser).Error; err != nil {
// 		response.WriteJson(w, http.StatusInternalServerError, map[string]string{
// 			"error": "Failed to create user",
// 		})
// 		return
// 	}

// 	// Respond success
// 	response.WriteJson(w, http.StatusOK, map[string]interface{}{
// 		"message": "User created successfully",
// 		 "user": modelUser,
// 	})
// }

// func LoginUser(w http.ResponseWriter,r *http.Request){
// 	//login user logic here
// 	var user types.User
// 	jwtSecret:=[]byte(os.Getenv("JWT_SECRET"))
// 	// Decode request body
// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		response.WriteJson(w, http.StatusBadRequest, map[string]string{
// 			"error": "Invalid request payload",
// 		})
// 		return
// 	}

// 	// Fetch user from database
// 	var modelUser model.User

// 	if err := database.DB.Where("email = ?", user.Email).First(&modelUser).Error; err != nil {
// 		response.WriteJson(w, http.StatusUnauthorized, map[string]string{
// 			"error": "User does not exit with this email",
// 		})
// 		return
// 	}

// 	// Compare passwords
// 	// bytePassword:=sha256.Sum256([]byte(modelUser.Password))
// 	// hashedPassword:=hex.EncodeToString(bytePassword[:])
// 	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 		if err != nil {
// 			response.WriteJson(w, http.StatusBadRequest, map[string]string{
// 			"error": "Invalid request payload",
// 		})
// 		return
// 		}
// 	if passwordHash != user.Password {
// 			response.WriteJson(w, http.StatusUnauthorized, map[string]string{
// 			"error": "Password is not correct!",
// 		})
// 		return
// 	}

// 	claims := jwt.MapClaims{
//         "sub": modelUser.Name,
// 		"ID":modelUser.ID,
//         "exp": time.Now().Add(time.Hour * 1).Unix(), // 1 hour expiry
//         "iat": time.Now().Unix(),                   // issued at
//     }
//     token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//     accessToken, err := token.SignedString(jwtSecret)
//     if err != nil {
//         response.WriteJson(w, http.StatusInternalServerError, map[string]string{
//             "error": "Failed to generate token",
//         })
//         return
//     }
// 	// Respond success
// 	response.WriteJson(w, http.StatusOK, map[string]interface{}{
// 		"message": "Login successful",
// 		 "user": modelUser,
// 		 "token":accessToken,
// 	})
// }



func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user types.User

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
		return
	}

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to hash password",
		})
		return
	}

	// Model to save
	modelUser := model.User{
		ID:       uuid.New().String(),
		Name:     user.Name,
		Email:    user.Email,
		Password: string(passwordHash),
	}

	// Save to DB
	if err := database.DB.Create(&modelUser).Error; err != nil {
		response.WriteJson(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user",
		})
		return
	}

	// Respond success
	response.WriteJson(w, http.StatusOK, map[string]interface{}{
		"message": "User created successfully",
		"user":    modelUser,
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user types.User
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
		return
	}

	// Fetch user from database
	var modelUser model.User
	if err := database.DB.Where("email = ?", user.Email).First(&modelUser).Error; err != nil {
		response.WriteJson(w, http.StatusUnauthorized, map[string]string{
			"error": "User does not exist with this email",
		})
		return
	}

	// Compare password correctly
	if err := bcrypt.CompareHashAndPassword([]byte(modelUser.Password), []byte(user.Password)); err != nil {
		response.WriteJson(w, http.StatusUnauthorized, map[string]string{
			"error": "Incorrect password",
		})
		return
	}

	// JWT Claims
	claims := jwt.MapClaims{
		"sub": modelUser.Email,
		"ID":  modelUser.ID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"iat": time.Now().Unix(),
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(jwtSecret)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
		return
	}

	// Success
	response.WriteJson(w, http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user":    modelUser,
		"token":   accessToken,
	})
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var resetPassword types.ResetPasswordRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&resetPassword); err != nil {
		response.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
		return
	}

	// Fetch only the Name field for given email
	var modelUser model.User
	if err := database.DB.
		Select("name").
		Where("email = ?", resetPassword.Email).
		First(&modelUser).Error; err != nil {

		response.WriteJson(w, http.StatusUnauthorized, map[string]string{
			"error": "User does not exist with this email",
		})
		return
	}

	// Hash the new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(resetPassword.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to hash password",
		})
		return
	}

	// ✅ Update password for this user
	if err := database.DB.
		Model(&model.User{}).
		Where("email = ?", resetPassword.Email).
		Update("password", string(passwordHash)).
		Error; err != nil {

		response.WriteJson(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to update password",
		})
		return
	}

	// Respond success (return only name)
	response.WriteJson(w, http.StatusOK, map[string]interface{}{
		"message": "Password reset successfully",
		"user": map[string]string{
			"name": modelUser.Name,
		},
		"user1":modelUser,
	})
}




