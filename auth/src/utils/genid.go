package utils

import (
	"github.com/google/uuid"
)

func GenID() string {
	// UUID を生成
	uidv4, _ := uuid.NewRandom()

	return uidv4.String()
}