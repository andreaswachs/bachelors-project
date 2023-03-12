package utils

import "github.com/google/uuid"

func RandomName() string {
	return uuid.New().String()
}
