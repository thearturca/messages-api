package messages

import (
	"context"
	"message-service/internal/db"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db      *pgxpool.Pool
	dialect *goqu.DialectWrapper
}

func NewService(db *pgxpool.Pool) *Service {
	dialect := goqu.Dialect("postgres")

	return &Service{
		db:      db,
		dialect: &dialect,
	}
}

func (s *Service) GetMessage(ctx context.Context, id string) (*db.Message, error) {
	query, args, err := s.dialect.Select(goqu.T("Messages").All()).Prepared(true).
		From(goqu.T("Messages")).Where(goqu.C("id").Eq(id)).
		ToSQL()

	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	message, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[db.Message])

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *Service) PostMessage(ctx context.Context, text string) (*db.Message, error) {
	query, args, err := s.dialect.Insert(goqu.T("Messages")).Prepared(true).
		Rows(db.Message{
			Text: text,
		}).
		Returning(goqu.T("Messages").All()).
		ToSQL()

	rows, err := s.db.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	message, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[db.Message])

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *Service) UpdateMessage(ctx context.Context, message *db.Message) (*db.Message, error) {
	query, args, err := s.dialect.Update(goqu.T("Messages")).
		Set(message).
		Where(goqu.C("id").Eq(message.Id)).
		Returning(goqu.T("Messages").All()).
		ToSQL()

	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	updatedMessage, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[db.Message])

	if err != nil {
		return nil, err
	}

	return &updatedMessage, nil
}
