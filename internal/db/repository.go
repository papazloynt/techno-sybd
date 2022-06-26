package db

import "github.com/jackc/pgx/v4/pgxpool"

type Repository struct {
	UserRepo    UserRepository
	ForumRepo   ForumRepository
	ThreadRepo  ThreadRepository
	PostRepo    PostRepository
	VoteRepo    VoteRepository
	ServiceRepo ServiceRepository
}

func NewRepository(db *pgxpool.Pool) (*Repository, error) {
	var err error
	repository := new(Repository)

	repository.UserRepo, err = NewUserRepository(db)
	if err != nil {
		return nil, err
	}

	repository.ForumRepo = NewForumRepository(db)

	repository.ThreadRepo, err = NewThreadRepository(db)
	if err != nil {
		return nil, err
	}

	repository.PostRepo, err = NewPostRepository(db)
	if err != nil {
		return nil, err
	}

	repository.VoteRepo, err = NewVoteRepository(db)
	if err != nil {
		return nil, err
	}

	repository.ServiceRepo = NewServiceRepository(db)

	return repository, nil
}
