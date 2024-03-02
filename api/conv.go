package api

import (
	"fmt"
	"rinha2024q1crebito"
	"strconv"
)

func parseInt(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf(
			"error parsing int: %w (%w)",
			err,
			rinha2024q1crebito.ErrInvalidParameter)
	}
	return v, nil
}
