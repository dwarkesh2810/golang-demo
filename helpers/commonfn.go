package helpers

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func SecondsToDayHourMinAndSeconds(seconds int) (int64, int64, int64, int64) {
	days := seconds / 86400
	hour := (seconds % 86400) / 3600
	minute := (seconds % 3600) / 60
	second := seconds % 60
	return int64(days), int64(hour), int64(minute), int64(second)
}

func HashData(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}
func VerifyHashedData(hashedString string, dataString string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(dataString))
	if err != nil {
		return errors.New("HASHED_ERROR")
	}
	return nil
}
