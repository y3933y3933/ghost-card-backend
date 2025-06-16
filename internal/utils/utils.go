package utils

import (
	"fmt"
	"strconv"
)

func ParseID(s string) (int64, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %v", err)
	}
	return id, nil
}
