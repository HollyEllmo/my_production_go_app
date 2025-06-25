package storage

import (
	"context"
	"errors"

	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/filter"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/sort"
	"github.com/HollyEllmo/my-first-go-project/pkg/client/postgresql"
	db "github.com/HollyEllmo/my-first-go-project/pkg/client/postgresql/model"
	"github.com/HollyEllmo/my-first-go-project/pkg/logging"
	sq "github.com/Masterminds/squirrel"
)

type ProductStorage struct {
	queryBuilder sq.StatementBuilderType
	client 	 PostgreSQLClient
}

func NewProductStorage(client PostgreSQLClient) *ProductStorage {
	return &ProductStorage{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client: client,
	}
}

const (
	scheme = "public"
	table = "product"
	tableScheme = scheme + "." + table
)


func (s *ProductStorage) All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]model.Product, error) {
	sortDB := db.NewSortOptions(sorting)
	filterDB := db.NewFilters(filtering)

	query := s.queryBuilder.Select("id").
		Column("name").
		Column("description").
		Column("image_id").
		Column("price").
		Column("currency_id").
		Column("rating").
		Column("created_at").
		Column("updated_at").
		From(tableScheme)

    query = filterDB.Enrich(query, "")
    query = sortDB.Sort(query, "")

	sql, args, err := query.ToSql()
	logger := logging.WithFields(ctx, map[string]interface{}{
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

func (s *ProductStorage) Create(ctx context.Context, dto *CreateProductStorageDTO) error {
	sql, args, buildError := s.queryBuilder.
	 Insert(tableScheme).
	 Columns(
		"id",
		"name",
		"description",
		"image_id",
		"price",
		"currency_id",
		"category_id",
		"specification",
		"rating",
		"created_at",
		"updated_at",
	 ).Values(
		dto.ID,
		dto.Name,
		dto.Description,
		dto.ImageID,
		dto.Price,
		dto.CurrencyID,
		dto.CategoryID,
		dto.Specification,
		dto.Rating,
		dto.CreatedAt,
		dto.UpdatedAt,
	 ).PlaceholderFormat(sq.Dollar).ToSql()

	 logger := logging.WithFields(ctx, map[string]interface{}{
		"sql":   sql,
		"table": tableScheme,
		"args":  args, 
	})
	if buildError != nil {
		buildError = db.ErrCreateQuery(buildError)
		logger.Error(buildError)
		return buildError
	}

	if exec, execErr := s.client.Exec(ctx, sql, args...); execErr != nil {
		execErr = db.ErrDoQuery(execErr)
		logger.Error(execErr)
		return execErr
	} else if exec.RowsAffected() == 0 || !exec.Insert() {
		execErr = db.ErrDoQuery(errors.New("product was not created. 0 raws were affected"))
		logger.Error(execErr)
		return execErr
	}

	return nil
}

func (s *ProductStorage) One(ctx context.Context, ID string) (*Product, error) {
	query := s.queryBuilder.Select("id").
		Column("name").
		Column("description").
		Column("image_id").
		Column("price").
		Column("currency_id").
		Column("rating").
		Column("category_id").
		Column("specification").
		Column("created_at").
		Column("updated_at").
		From(tableScheme).
		Where(sq.Eq{"id": ID})

	sql, args, err := query.ToSql()
	logger := logging.WithFields(ctx, map[string]interface{}{
		"sql":   sql,
		"table": table,
		"args":  args,
		"id":    ID,
	})
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return nil, err
	}

	logger.Trace("do query")
	var p Product
	if err = s.client.QueryRow(ctx, sql, args...).Scan(
		&p.ID, &p.Name, &p.Description, &p.ImageID, &p.Price, 
		&p.CurrencyID, &p.Rating, &p.CategoryID, &p.Specification, 
		&p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		err = db.ErrScan(postgresql.ParsePgError(err))
		logger.Error(err)
		return nil, err
	}

	return &p, nil
}