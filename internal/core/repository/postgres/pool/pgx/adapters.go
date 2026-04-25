package core_pgx_pool

import (
	"errors"
	"fmt"

	core_postgres_pool "github.com/artlink52/go-todoapp/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgxRows struct {
	pgx.Rows
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return mapErrors(err)
	}
	return nil
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

func mapErrors(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23503" {
			return fmt.Errorf(
				"%v: %w",
				err,
				core_postgres_pool.ErrViolatesForeignKey,
			)
		}
	}
	return fmt.Errorf(
		"%v: %w",
		err,
		core_postgres_pool.ErrUnknown,
	)

}
