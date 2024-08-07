package queryservice

import (
	"context"
	"datapoint/internal/model/dbmodel"
	"datapoint/internal/model/querymodel"
	"fmt"
	"go.uber.org/zap"
)

type DBService interface {
	GetByID(id string) (*dbmodel.DB, error)
}

type service struct {
	dbService DBService
}

func (s *service) Execute(ctx context.Context, info querymodel.Info, id string) (querymodel.QueryResult, error) {
	zap.S().Info("попытка выполнить запрос")

	db, err := s.dbService.GetByID(id)
	if err != nil {
		return querymodel.QueryResult{}, err
	}

	q := querymodel.New(info, db.B())

	var result querymodel.QueryResult
	if result, err = q.Execute(ctx, db); err != nil {
		err = fmt.Errorf("не удалось выполнить запрос: %s", err)
		zap.S().Error(err)
		return querymodel.QueryResult{}, err
	}

	zap.S().Info("запрос выполнен успешно")
	return result, nil
}

func New(dbService DBService) *service {
	return &service{dbService: dbService}
}
