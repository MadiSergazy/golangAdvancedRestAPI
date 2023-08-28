package author

import (
	"context"
	"errors"
	"fmt"
	"strings"

	// "github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgconn"

	"mado/internal/author"
	"mado/pkg/client/postgresql"
	"mado/pkg/logging"
)

type repository struct {
	client postgresql.Client

	logger *logging.Logger
}

func formatQuery(query string) string {
	return strings.ReplaceAll(strings.ReplaceAll(query, "\t", ""), "\n", "")
}

// Create implements author.Repository.
func (r *repository) Create(ctx context.Context, author *author.Author) error {
	q := `INSERT INTO public.author (name) VALUES ($1) RETURNING id`

	r.logger.Trace(fmt.Sprintf("SQL Query: ", formatQuery(q)))
	var pgErr *pgconn.PgError

	row := r.client.QueryRow(ctx, q, author.Name)
	if err := row.Scan(&author.ID); err != nil {

		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}

	return nil

}

// Delete implements author.Repository.
func (r *repository) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// FindAll implements author.Repository.
func (r *repository) FindAll(ctx context.Context) (u []author.Author, err error) {
	q := `SELECT id, name FROM public.author`
	r.logger.Trace(fmt.Sprintf("SQL Query: ", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	authors := make([]author.Author, 0)
	for rows.Next() {
		var auth author.Author
		if err = rows.Scan(&auth.ID, auth.Name); err != nil {
			return nil, err
		}

		authors = append(authors, auth)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return authors, nil
}

// FindOne implements author.Repository.
func (r *repository) FindOne(ctx context.Context, id string) (author.Author, error) {
	q := `
SELECT id, name 
FROM public.author 
WHERE id = $1`
	r.logger.Trace(fmt.Sprintf("SQL Query: ", formatQuery(q)))

	var auth author.Author
	row := r.client.QueryRow(ctx, q, id)
	if err := row.Scan(&auth.ID, &auth.Name); err != nil {
		return author.Author{}, nil
	}

	return auth, nil
}

// Update implements author.Repository.
func (r *repository) Update(ctx context.Context, user author.Author) error {
	panic("unimplemented")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) author.Repository {

	return &repository{
		client: client,
		logger: logger,
	}
}
