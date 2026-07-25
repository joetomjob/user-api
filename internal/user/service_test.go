package user

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
)

type FakeRepo struct {
	created User
}

func (f *FakeRepo) Create(ctx context.Context, user User) (User, error) {
	user.Id = 1
	f.created = user
	return f.created, nil
}
func (f *FakeRepo) GetById(ctx context.Context, id int) (User, error) {
	if id == 100 {
		return User{}, pgx.ErrNoRows
	}
	return User{Id: id, Name: "Joe", Email: "joetomjob@gmail.com", Age: 36}, nil
}
func (f *FakeRepo) Update(ctx context.Context, user User, id int) (User, error) {
	user.Id = id
	f.created = user
	return f.created, nil
}
func (f *FakeRepo) Delete(ctx context.Context, id int) error {
	return nil
}

func TestService(t *testing.T) {
	t.Run("CRUD Success", func(t *testing.T) {
		s := NewService(&FakeRepo{})

		user := User{Name: "Joe", Email: "joetomjob@gmail.com", Age: 36}

		gotUser, err := s.Create(t.Context(), user)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if gotUser.Id != 1 {
			t.Fatal("Invalid User id")
		}

		newId := 2
		gotUser, err = s.Update(t.Context(), user, 2)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if gotUser.Id != newId {
			t.Fatal("Invalid User id")
		}

		newId = 10
		gotUser, err = s.GetById(t.Context(), newId)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if gotUser.Id != newId {
			t.Fatal("Invalid User id")
		}

		err = s.Delete(t.Context(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

	})
}
