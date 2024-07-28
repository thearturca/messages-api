package statistics

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatisticsDto struct {
	TotalMessages     int `db:"total_messages" json:"totalMessages"`
	ProcessedMessages int `db:"processed_messages" json:"processedMessages"`
}

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

func (s *Service) GetStatistics(ctx context.Context, from *time.Time, to *time.Time) (*StatisticsDto, error) {
	query, args, err := s.dialect.Select(
		goqu.COUNT(goqu.T("M").Col("id")).As("total_messages"),
		goqu.L("COUNT(\"M\".id) FILTER (WHERE \"M\".is_processed is true) as processed_messages"),
	).
		From(goqu.T("Messages").As("M")).
		Where(
			goqu.T("M").Col("created_at").Between(goqu.Range(
				goqu.COALESCE(from, goqu.L("'-infinity'::timestamp")),
				goqu.COALESCE(to, goqu.L("'infinity'::timestamp")),
			)),
		).
		ToSQL()

	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	stats, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[StatisticsDto])

	if err != nil {
		return nil, err
	}

	return &stats, nil
}
