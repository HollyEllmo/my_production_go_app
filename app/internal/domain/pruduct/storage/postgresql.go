package storage

import (
	"context"

	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/HollyEllmo/my-first-go-project/pkg/client/postgresql"
	db "github.com/HollyEllmo/my-first-go-project/pkg/client/postgresql/model"
	"github.com/HollyEllmo/my-first-go-project/pkg/logging"
	sq "github.com/Masterminds/squirrel"
)

type ProductStorage struct {
	queryBuilder sq.StatementBuilderType
	client 	 PostgreSQLClient
}

func NewProductStorage(client PostgreSQLClient) ProductStorage {
	return ProductStorage{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client: client,
	}
}

const (
	scheme = "public"
	table = "product"
)

// TODO вынести этот метод в отдельный пакет, например, в pkg/logging

func (s *ProductStorage) All(ctx context.Context) ([]model.Product, error) {
	query := s.queryBuilder.Select("id").
		Column("name").
		Column("description").
		Column("image_id").
		Column("price").
		Column("currency_id").
		Column("rating").
		Column("created_at").
		Column("updated_at").
		From(scheme + "." + table)

	// TODO Реализовать filtering and sorting по полям

	sql, args, err := query.ToSql()
	logger := logging.GetLogger(ctx).WithFields(map[string]interface{}{
		"sql":   sql,
		"table": table,
		"args":  args, 
	})
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return nil, err
	}

	logger.Trace("do query")
	rows, err := s.client.Query(ctx, sql, args...)
	if err != nil {
		err = db.ErrDoQuery(err)
		logger.Error(err)
		return nil, err
	}

	defer rows.Close()

	list := make([]model.Product, 0)

	for rows.Next() {
		p := model.Product{}
		if err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImageID, &p.Price, &p.CurrencyID, &p.Rating, &p.CreatedAt, &p.UpdatedAt); 
		err != nil {
			err = db.ErrScan(postgresql.ParsePgError(err))
			logger.Error(err)
			return nil, err
		}

		list = append(list, p)
	}
	

	return list, nil
}