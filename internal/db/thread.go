package db

import (
	"SYBD/internal/model/core"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	// INSERT
	qCreateThread = "INSERT INTO \"thread\" (title, author, forum, message, slug, created) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, title, author, forum, message, votes, slug, created;"

	//UPDATE
	qUpdateThread = "UPDATE \"thread\" SET title = $2, message = $3 WHERE id = $1 RETURNING id, title, author, forum, message, votes, slug, created;"

	// SELECT
	qGetThreadBySlug = "SELECT id, title, author, forum, message, votes, slug, created FROM \"thread\" WHERE slug = $1;"
	qGetThreadByID   = "SELECT id, title, author, forum, message, votes, slug, created FROM \"thread\" WHERE id = $1;"
)

type ThreadRepository interface {
	CreateThread(ctx context.Context, thread *core.Thread) (*core.Thread, error)
	UpdateThread(ctx context.Context, id int64, title string, message string) (*core.Thread, error)
	GetThread(ctx context.Context, slug string) (*core.Thread, error)
	GetThreadByID(ctx context.Context, id int64) (*core.Thread, error)
}

type threadRepositoryImpl struct {
	dbConn *pgxpool.Pool
}

func (repo *threadRepositoryImpl) CreateThread(ctx context.Context, thread *core.Thread) (*core.Thread, error) {
	t := &core.Thread{}
	err := repo.dbConn.QueryRow(ctx,
		qCreateThread,
		thread.Title,
		thread.Author,
		thread.Forum,
		thread.Message,
		thread.Slug,
		thread.Created).
		Scan(&t.ID,
			&t.Title,
			&t.Author,
			&t.Forum,
			&t.Message,
			&t.Votes,
			&t.Slug,
			&t.Created)
	return t, err
}

func (repo *threadRepositoryImpl) GetThread(ctx context.Context, slug string) (*core.Thread, error) {
	t := &core.Thread{}
	err := repo.dbConn.QueryRow(ctx,
		qGetThreadBySlug,
		slug).
		Scan(&t.ID,
			&t.Title,
			&t.Author,
			&t.Forum,
			&t.Message,
			&t.Votes,
			&t.Slug,
			&t.Created)
	return t, wrapErr(err)
}

func (repo *threadRepositoryImpl) GetThreadByID(ctx context.Context, id int64) (*core.Thread, error) {
	t := &core.Thread{}
	err := repo.dbConn.QueryRow(ctx,
		qGetThreadByID,
		id).
		Scan(&t.ID,
			&t.Title,
			&t.Author,
			&t.Forum,
			&t.Message,
			&t.Votes,
			&t.Slug,
			&t.Created)
	return t, wrapErr(err)
}

func (repo *threadRepositoryImpl) UpdateThread(ctx context.Context, id int64, title string, message string) (*core.Thread, error) {
	t := &core.Thread{}
	err := repo.dbConn.QueryRow(ctx,
		qUpdateThread,
		id,
		title,
		message).
		Scan(&t.ID,
			&t.Title,
			&t.Author,
			&t.Forum,
			&t.Message,
			&t.Votes,
			&t.Slug,
			&t.Created)
	return t, wrapErr(err)
}

func NewThreadRepository(dbConn *pgxpool.Pool) (*threadRepositoryImpl, error) {
	return &threadRepositoryImpl{dbConn: dbConn}, nil
}
