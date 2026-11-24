package repository

import sqlcdb "threadly-api-v1/internal/infrastructure/database/sqlc"

type UserRepository struct {
	queries *sqlcdb.Queries
}

func NewUserRepository(queries *sqlcdb.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) GetByID(id string) (interface{}, error) {
	return nil, nil
}
