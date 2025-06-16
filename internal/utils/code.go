package utils

import (
	"context"
	"database/sql"
	"errors"

	"math/rand"

	"github.com/y3933y3933/joker/internal/database"
)

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var ErrGenerateCode = errors.New("failed to generate unique game code")

func GenerateUniqueGameCode(ctx context.Context, q *database.Queries, length, maxRetries int) (string, error) {
	for i := 0; i < maxRetries; i++ {
		code := RandomCode(length)
		_, err := q.GetGameByCode(ctx, code)

		if errors.Is(err, sql.ErrNoRows) {
			return code, nil
		}

		if err != nil {
			return "", err
		}
	}
	return "", ErrGenerateCode
}
