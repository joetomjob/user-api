package user

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func TestRepository(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Skipf("Skipping integration test: failed to fetch env file: %s", err)
	}

	url := os.Getenv("DATABASE_URL")
	if len(strings.TrimSpace(url)) == 0 {
		t.Skip("Skipping integration test: database url is empty")
	}

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Skipf("Skipping integration test: Config createion failed: %s", err)
	}
	config.MaxConnIdleTime = time.Minute * 30
	config.MaxConns = 20
	config.MinConns = 10

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Skipf("Skipping integration test: Pool creattion failed: %s", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		t.Skipf("Skipping integration test: Ping to db failed: %s", err)
	}

	repo := NewRepo(pool)
	user := User{Name: "test", Email: "test@gmail.com", Age: 36}
	t.Run("Create test", func(t *testing.T) {
		got, err := repo.Create(ctx, user)

		if err != nil {
			t.Fatalf("user creation failed: %s", err)
		}
		if got.Id < 0 || user.Name != got.Name || user.Email != got.Email || user.Age != got.Age {
			t.Fatalf("invalid data")
		}
		user.Id = got.Id
	})

	t.Run("GetById test", func(t *testing.T) {
		got, err := repo.GetById(ctx, user.Id)

		if err != nil {
			t.Fatalf("user fetch failed: %s", err)
		}
		if user.Id != got.Id {
			t.Fatalf("invalid data")
		}
	})

	t.Run("Update test", func(t *testing.T) {
		user.Name = "testtest"
		got, err := repo.Update(ctx, user, user.Id)

		if err != nil {
			t.Fatalf("user update failed: %s", err)
		}
		if got.Id < 0 || user.Name != got.Name || user.Email != got.Email || user.Age != got.Age {
			t.Fatalf("invalid data")
		}
	})

	t.Run("Delete test", func(t *testing.T) {
		err := repo.Delete(ctx, user.Id)

		if err != nil {
			t.Fatalf("user delete failed: %s", err)
		}
	})
}
