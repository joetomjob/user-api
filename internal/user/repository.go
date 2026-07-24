package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("resource not found")

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool}
}

func (repo *Repo) Create(ctx context.Context, user User) (User, error) {
	var id int
	query := "INSERT INTO users (name, email, age) values ($1, $2, $3) RETURNING id"
	err := repo.pool.QueryRow(ctx, query, user.Name, user.Email, user.Age).Scan(&id)
	if err != nil {
		return User{}, err
	}

	user.Id = id
	return user, nil

}

func (repo *Repo) GetById(ctx context.Context, id int) (User, error) {
	var user User

	query := "SELECT id, name, email, age from users where id = $1;"

	err := repo.pool.QueryRow(ctx, query, id).Scan(&user.Id, &user.Name, &user.Email, &user.Age)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (repo *Repo) Update(ctx context.Context, user User, id int) (User, error) {
	query := "UPDATE users SET name = $1, email = $2, age = $3 WHERE id = $4"
	ct, err := repo.pool.Exec(ctx, query, user.Name, user.Email, user.Age, id)
	if err != nil {
		return User{}, err
	}

	if ct.RowsAffected() == 0 {
		return User{}, ErrNotFound
	}

	user.Id = id
	return user, nil

}

func (repo *Repo) Delete(ctx context.Context, id int) error {
	query := "DELETE from users where id = $1;"

	ct, err := repo.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
