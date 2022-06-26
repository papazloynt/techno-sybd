package db

import (
	"SYBD/internal/model/core"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	// INSERT
	qCreateVote = "INSERT INTO \"vote\" (nickname, thread, voice) VALUES ($1, $2, $3);"

	// SELECT
	qExists = "SELECT voice from \"vote\" WHERE nickname = $1 AND thread = $2;"
	qUpdate = "UPDATE \"vote\" SET voice = $3 WHERE thread = $1 AND nickname = $2 AND voice != $3;"
)

type VoteRepository interface {
	CreateVote(ctx context.Context, vote *core.Vote) error
	VoteExists(ctx context.Context, nickname string, threadID int64) (bool, error)
	UpdateVote(ctx context.Context, threadID int64, nickname string, voice int64) (bool, error)
}

type voteRepositoryImpl struct {
	db *pgxpool.Pool
}

func (repo *voteRepositoryImpl) CreateVote(ctx context.Context, vote *core.Vote) error {
	_, err := repo.db.Exec(ctx,
		qCreateVote,
		vote.Nickname,
		vote.ThreadID,
		vote.Voice)
	return wrapErr(err)
}

func (repo *voteRepositoryImpl) VoteExists(ctx context.Context, nickname string, threadID int64) (bool, error) {
	voice := 0
	err := repo.db.QueryRow(ctx,
		qExists,
		nickname,
		threadID).Scan(&voice)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (repo *voteRepositoryImpl) UpdateVote(ctx context.Context, threadID int64, nickname string, voice int64) (bool, error) {
	res, err := repo.db.Exec(ctx,
		qUpdate,
		threadID,
		nickname,
		voice)
	if err != nil {
		return false, err
	}
	return res.RowsAffected() == 1, nil
}

func NewVoteRepository(db *pgxpool.Pool) (*voteRepositoryImpl, error) {
	return &voteRepositoryImpl{db: db}, nil
}
