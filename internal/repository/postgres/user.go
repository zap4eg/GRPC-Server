package postgres

import (
	"GRPCProject/internal/core"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (repository *UserRepository) GetById(ctx context.Context, id string) (*core.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userChannel := make(chan *core.User)

	var err error
	go func() {
		err = repository.retrieveUser(ctx, id, userChannel)
	}()
	if err != nil {
		return nil, err
	}

	select {
	case user := <-userChannel:
		return user, nil
	case <-ctxTimeout.Done():
		{
			return nil, ctxTimeout.Err()
		}
	}
}

func (repository *UserRepository) retrieveUser(ctx context.Context, id string, channel chan<- *core.User) error {
	userId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	user := &core.User{}

	err = repository.pool.QueryRow(ctx,
		"SELECT id, first_name, last_name FROM users WHERE id = $1", userId).
		Scan(&user.ID, &user.FirstName, &user.LastName)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	channel <- user
	return nil
}
