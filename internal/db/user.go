package db

import (
	"SYBD/internal/model/core"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	// INSERT Query
	qCreateUser = "INSERT INTO \"user\" (nickname, fullname, about, email) VALUES ($1, $2, $3, $4);"

	// SELECT Query
	qGetUserByNickname = "SELECT nickname, fullname, about, email FROM \"user\" WHERE nickname = $1;"
	qGetUserByEmail    = "SELECT nickname, fullname, about, email FROM \"user\" WHERE email = $1;"
	qGetSimilaryUsers  = "SELECT nickname, fullname, about, email FROM \"user\" WHERE email = $1 OR nickname = $2;"
	qUpdateUser        = "UPDATE \"user\" SET fullname = COALESCE(NULLIF(TRIM($1), ''), fullname), about = COALESCE(NULLIF(TRIM($2), ''), about), email = COALESCE(NULLIF(TRIM($3), ''), email) WHERE nickname = $4 RETURNING fullname, about, email;"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *core.User) error
	GetUserByNickname(ctx context.Context, nickname string) (*core.User, error)
	GetUserByEmail(ctx context.Context, email string) (*core.User, error)
	GetSimilaryUsers(ctx context.Context, email string, nickname string) ([]core.User, error)
	UpdateUser(ctx context.Context, user *core.User) (*core.User, error)
}

type userRepositoryImpl struct {
	db *pgxpool.Pool
}

func (repo *userRepositoryImpl) CreateUser(ctx context.Context, user *core.User) error {
	_, err := repo.db.Exec(ctx,
		qCreateUser,
		user.Nickname,
		user.FullName,
		user.About,
		user.Email)
	return wrapErr(err)
}

func (repo *userRepositoryImpl) GetUserByNickname(ctx context.Context, nickname string) (*core.User, error) {
	user := &core.User{}
	err := repo.db.QueryRow(ctx,
		qGetUserByNickname,
		nickname).
		Scan(&user.Nickname,
			&user.FullName,
			&user.About,
			&user.Email)
	return user, wrapErr(err)
}

func (repo *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*core.User, error) {
	user := &core.User{}
	err := repo.db.QueryRow(ctx,
		qGetUserByEmail,
		email).
		Scan(&user.Nickname,
			&user.FullName,
			&user.About,
			&user.Email)
	return user, wrapErr(err)
}

func (repo *userRepositoryImpl) GetSimilaryUsers(ctx context.Context, email string, nickname string) ([]core.User, error) {
	rows, err := repo.db.Query(ctx,
		qGetSimilaryUsers,
		email,
		nickname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []core.User
	for rows.Next() {
		user := &core.User{}
		if err := rows.Scan(&user.Nickname,
			&user.FullName,
			&user.About,
			&user.Email); err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, nil
}

func (repo *userRepositoryImpl) UpdateUser(ctx context.Context, user *core.User) (*core.User, error) {
	updatedUser := &core.User{Nickname: user.Nickname}
	if err := repo.db.QueryRow(ctx, qUpdateUser,
		user.FullName, user.About,
		user.Email, user.Nickname).
		Scan(&updatedUser.FullName,
			&updatedUser.About,
			&updatedUser.Email); err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func NewUserRepository(db *pgxpool.Pool) (*userRepositoryImpl, error) {
	return &userRepositoryImpl{db: db}, nil
}
