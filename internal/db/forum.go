package db

import (
	"SYBD/internal/model/core"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/jackc/pgx/v4"
)

const (
	// INSERT
	qCreateForum = `INSERT INTO "forum" (title, "user", slug) VALUES ($1, $2, $3);`

	// SELECT
	qGetForumBySlug = `SELECT title, "user", slug, posts, threads FROM "forum" WHERE slug = $1;`
)

type ForumRepository interface {
	CreateForum(ctx context.Context, forum *core.Forum) error
	GetForumBySlug(ctx context.Context, slug string) (*core.Forum, error)
	GetUsersFromForum(ctx context.Context, slug string, limit int64, since string, desc bool) ([]*core.User, error)
	GetThreadsFromForum(ctx context.Context, slug string, limit int64, since string, desc bool) ([]*core.Thread, error)
}

type forumRepositoryImpl struct {
	db *pgxpool.Pool
}

func (repo *forumRepositoryImpl) CreateForum(ctx context.Context, forum *core.Forum) error {
	_, err := repo.db.Exec(ctx,
		qCreateForum,
		&forum.Title,
		&forum.User,
		&forum.Slug)
	return err
}

func (repo *forumRepositoryImpl) GetForumBySlug(ctx context.Context, slug string) (*core.Forum, error) {
	forum := &core.Forum{}
	err := repo.db.QueryRow(ctx,
		qGetForumBySlug,
		slug).
		Scan(&forum.Title,
			&forum.User,
			&forum.Slug,
			&forum.Posts,
			&forum.Threads)
	return forum, wrapErr(err)
}

func (repo *forumRepositoryImpl) GetUsersFromForum(ctx context.Context, slug string, limit int64, since string, desc bool) ([]*core.User, error) {
	query := constructGetForumUsersQuery(limit, since, desc)
	rows, err := repo.db.Query(ctx, query, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*core.User, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		u := &core.User{}
		if err := rows.Scan(&u.Nickname, &u.FullName, &u.About, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

const (
	qTemplate = "SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created FROM \"thread\" as t LEFT JOIN \"forum\" f ON t.forum = f.slug WHERE f.slug = $1 "
)

func (repo *forumRepositoryImpl) GetThreadsFromForum(ctx context.Context, slug string, limit int64, since string, desc bool) ([]*core.Thread, error) {
	var rows pgx.Rows
	var err error
	// Create query with conditions
	query := qTemplate

	queryOrderBy := "ORDER BY t.created "
	if desc {
		queryOrderBy += "DESC "
	}
	if limit > 0 {
		queryOrderBy += fmt.Sprintf("LIMIT %d ", limit)
	}

	if since != "" {
		querySince := "AND t.created >= $2 "
		if since != "" && desc {
			querySince = "AND t.created <= $2 "
		} else if since != "" && !desc {
			querySince = "AND t.created >= $2 "
		}

		query += querySince + queryOrderBy
		rows, err = repo.db.Query(ctx, query, slug, since)
	} else {
		query += queryOrderBy
		rows, err = repo.db.Query(ctx, query, slug)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threads := make([]*core.Thread, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		t := &core.Thread{}
		if err := rows.Scan(&t.ID, &t.Title,
			&t.Author, &t.Forum,
			&t.Message, &t.Votes,
			&t.Slug, &t.Created); err != nil {
			return nil, err
		}
		threads = append(threads, t)
	}

	return threads, nil
}

func NewForumRepository(db *pgxpool.Pool) *forumRepositoryImpl {
	return &forumRepositoryImpl{db: db}
}

func constructGetForumUsersQuery(limit int64, since string, desc bool) string {
	query := "SELECT u.nickname, u.fullname, u.about, u.email from \"forum_user\" u where u.forum = $1 "

	if len(since) > 0 {
		if desc {
			query += fmt.Sprintf("and u.nickname < '%s' ", since)
		} else {
			query += fmt.Sprintf("and u.nickname > '%s' ", since)
		}
	}

	query += "ORDER BY u.nickname "
	if desc {
		query += "DESC "
	}
	query += fmt.Sprintf("LIMIT %d ", limit)

	return query
}
