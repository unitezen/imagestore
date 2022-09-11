package utilities

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/unitezen/imagestore/core"
	"github.com/unitezen/imagestore/models"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePlaintextAndHash(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func GenerateRandomAPIKey() (string, error) {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, 32)
	for i := 0; i < 32; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}
	return string(ret), nil
}

func CreateUserAPIKey(user *models.User) (*models.UserSession, error) {

	// Delete any existing API key
	DeleteUserAPIKey(user)

	// Create a new key for the user
	apiKey, err := GenerateRandomAPIKey()
	if err != nil {
		return nil, err
	}

	// Persist the API key to the database
	userSession := &models.UserSession{ApiKey: apiKey, UserId: user.ID}
	result := core.Database.Create(userSession)
	if result.Error != nil {
		return nil, err
	}

	return userSession, nil
}

func DeleteUserAPIKey(user *models.User) (bool, error) {
	result := core.Database.Where("user_id = ?", user.ID).Delete(&models.UserSession{})
	if result.RowsAffected == 0 {
		return false, errors.New("No session with matching user ID found for deletion")
	}
	return true, nil
}
