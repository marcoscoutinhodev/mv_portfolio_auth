package user

import (
	"context"
	"database/sql"

	"github.com/marcoscoutinhodev/ms_auth/internal/entity"
	"github.com/marcoscoutinhodev/ms_auth/internal/infra/db/postgres"
)

type Repository struct {
	db      *sql.DB
	queries *postgres.Queries
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db:      db,
		queries: postgres.New(db),
	}
}

func (r Repository) Find(ctx context.Context, id string) (*entity.User, error) {
	user, err := r.queries.Find(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return entity.NewUser(user.ID, user.Name, user.Email, user.Password), nil
}

func (r Repository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	u, err := r.queries.FindByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	user := entity.NewUser(u.ID, u.Name, u.Email, u.Password)
	user.ConfirmedEmail = u.ConfirmedEmail.Bool
	return user, nil
}

func (r Repository) Store(ctx context.Context, user *entity.User, fn func() error) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return err
	}

	qtx := r.queries.WithTx(tx)
	if err := qtx.Store(ctx, postgres.StoreParams{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}); err != nil {
		return err
	}

	if err := fn(); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r Repository) Update(ctx context.Context, user *entity.User) error {
	if err := r.queries.Update(ctx, postgres.UpdateParams{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}); err != nil {
		return err
	}

	return nil
}

func (r Repository) ConfirmEmail(ctx context.Context, userID string) error {
	if err := r.queries.ConfirmEmail(ctx, userID); err != nil {
		return err
	}

	return nil
}
