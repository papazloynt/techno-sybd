package db

import (
	"SYBD/internal/model/core"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	// TRUNCATE
	qDeleteTables = "TRUNCATE TABLE \"user\", \"forum\", \"thread\", \"post\", \"forum_user\", \"vote\" CASCADE;"

	// SELECT
	qCountAll = "SELECT (SELECT count(*) FROM \"user\") AS user, (SELECT count(*) FROM \"forum\") AS forum, (SELECT count(*) FROM \"thread\") AS thread, (SELECT count(*) FROM \"post\") AS post;"
)

type ServiceRepository interface {
	Status(ctx context.Context) (*core.Service, error)
	Delete(ctx context.Context) error
}

type serviceRepositoryImpl struct {
	db *pgxpool.Pool
}

func (repo *serviceRepositoryImpl) Status(ctx context.Context) (*core.Service, error) {
	res := &core.Service{}
	err := repo.db.QueryRow(ctx, qCountAll).
		Scan(&res.User, &res.Forum,
			&res.Thread, &res.Post)
	return res, err
}

func (repo *serviceRepositoryImpl) Delete(ctx context.Context) error {
	_, err := repo.db.Exec(ctx, qDeleteTables)
	return err
}

func NewServiceRepository(db *pgxpool.Pool) *serviceRepositoryImpl {
	return &serviceRepositoryImpl{db: db}
}
