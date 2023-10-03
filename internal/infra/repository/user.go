package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/marcoscoutinhodev/mv_chat/internal/entity"
	"github.com/marcoscoutinhodev/mv_chat/internal/infra/repository/postgres"
)

type UserRepository struct {
	db        *sql.DB
	txOptions *sql.TxOptions
	queries   *postgres.Queries
}

func (u UserRepository) Find(ctx context.Context, id string) (*entity.User, error) {
	user, err := u.queries.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	return entity.NewUser(user.ID, user.Name, user.Email, user.Password), nil
}

func (u UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.queries.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return entity.NewUser(user.ID, user.Name, user.Email, user.Password), nil
}

func (u UserRepository) Store(ctx context.Context, user *entity.User, fn func() error) error {
	tx, err := u.db.BeginTx(ctx, u.txOptions)
	if err != nil {
		return err
	}

	qtx := u.queries.WithTx(tx)
	if err := qtx.Store(ctx, postgres.StoreParams{
		ID:       uuid.NewString(),
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