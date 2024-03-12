package database

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
)

const (
	CodeErrForeign    = "23503"
	CodeErrExists     = "23505"
	CodeErrValidation = "23514"
)

type ctxTXClient struct{}

func (db Database) Exec(ctx context.Context, sql string, args ...any) error {
	var client interface {
		Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	} = db.client
	if tx, ok := ctx.Value(ctxTXClient{}).(pgx.Tx); ok {
		client = tx
	}

	_, err := client.Exec(ctx, sql, args...)

	return databaseError(err)
}

func (db Database) QueryAll(ctx context.Context, dst any, sql string, args ...any) error {
	var client interface {
		Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	} = db.client
	if tx, ok := ctx.Value(ctxTXClient{}).(pgx.Tx); ok {
		client = tx
	}

	rows, err := client.Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err := pgxscan.ScanAll(dst, rows); err != nil {
		if pgxscan.NotFound(err) {
			return model.NotFoundError(err)
		}

		return databaseError(err)
	}

	return nil
}

func (db Database) QueryOne(ctx context.Context, dst any, sql string, args ...any) error {
	var client interface {
		Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	} = db.client
	if tx, ok := ctx.Value(ctxTXClient{}).(pgx.Tx); ok {
		client = tx
	}

	rows, err := client.Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err := pgxscan.ScanOne(dst, rows); err != nil {
		if pgxscan.NotFound(err) {
			return model.NotFoundError(err)
		}

		return databaseError(err)
	}

	return nil
}

func databaseError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case CodeErrValidation:
			return model.WrappedError(model.DatabaseError{
				Name: pgErr.ConstraintName,
				Code: pgErr.Code,
				Err:  err,
			}, "database validation failed")
		case CodeErrForeign, CodeErrExists:
			return model.DatabaseError{Name: pgErr.ConstraintName, Code: pgErr.Code, Err: err}
		default:
			return model.DatabaseError{Code: pgErr.Code, Err: err}
		}
	}
	return err
}

func (db Database) ContextTransactionBegin(ctx context.Context) (context.Context, error) {
	if tx, ok := ctx.Value(ctxTXClient{}).(pgx.Tx); ok {
		tx, err := tx.Begin(ctx)
		if err != nil {
			return nil, err
		}
		return context.WithValue(ctx, ctxTXClient{}, tx), nil
	}

	tx, err := db.client.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, ctxTXClient{}, tx), nil
}

func (db Database) ContextTransactionCommit(ctx context.Context) error {
	if tx, ok := ctx.Value(ctxTXClient{}).(pgx.Tx); ok {
		return tx.Commit(ctx)
	}
	return errors.New("no transaction to commit")
}

func (db Database) ContextTransactionRollback(ctx context.Context) error {
	if tx, ok := ctx.Value(ctxTXClient{}).(pgx.Tx); ok {
		return tx.Rollback(ctx)
	}
	return errors.New("no transaction to rollback")
}
