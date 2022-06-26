package db

import (
	"SYBD/internal/constants"
	"errors"
	"github.com/jackc/pgx/v4"
)

func wrapErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return constants.ErrDBNotFound
	}

	return err
}
