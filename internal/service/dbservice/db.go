package dbservice

import (
	"context"
	"datapoint/internal/model/dbmodel"
	"datapoint/pkg/database"
	"errors"
	"fmt"
	"go.uber.org/zap"
)

type DBRepo interface {
	GetList(ctx context.Context) ([]*dbmodel.DB, error)
	Add(ctx context.Context, db dbmodel.DB) error
	Edit(ctx context.Context, db dbmodel.DB) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	r      DBRepo
	tx     database.TxManager
	dbList map[string]*dbmodel.DB
}

func (s *service) GetList() []*dbmodel.DB {
	list := make([]*dbmodel.DB, 0, len(s.dbList))
	for _, d := range s.dbList {
		list = append(list, d)
	}
	return list
}

func (s *service) GetByID(id string) (*dbmodel.DB, error) {
	zap.S().Info("попытка получить базу данных по идентификатору",
		zap.String("id", id))

	db, ok := s.dbList[id]
	if !ok {
		err := errors.New("базы данных не существует")
		zap.S().Error(err, zap.String("id", id))
		return nil, err
	}

	zap.S().Info("база данных успешно получена", zap.String("id", id))
	return db, nil
}

func (s *service) Add(ctx context.Context, info dbmodel.Info) (string, error) {
	zap.S().Info("попытка добавить базу данных")

	db, err := dbmodel.New(info)
	if err != nil {
		err = fmt.Errorf("не удалось подключиться к базе данных: %s", err)
		zap.S().Error(err)
		return "", err
	}

	if err = s.r.Add(ctx, *db); err != nil {
		db.Close()
		err = fmt.Errorf("не удалось сохранить базу данных: %s", err)
		zap.S().Error(err)
		return "", err
	}

	s.dbList[db.ID] = db

	zap.S().Info("база данных успешно добавлена")
	return db.ID, nil
}

func (s *service) Edit(ctx context.Context, info dbmodel.Info, id string) error {
	zap.S().Info("попытка отредактировать базу данных", zap.String("id", id))

	db, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if err = s.tx.ReadCommitted(ctx, func(ctx context.Context) error {
		err := s.r.Edit(ctx, *db)
		if err != nil {
			err = fmt.Errorf("не удалось отредактировать базу данных: %s", err)
			zap.S().Error(err)
			return err
		}

		if err = db.SetInfo(info); err != nil {
			err = fmt.Errorf("не удалось изменить параметры подключения: %s", err)
			zap.S().Error(err)
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	zap.S().Info("база данных успешно отредактирована", zap.String("id", id))
	return nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	zap.S().Info("попытка удалить базу данных", zap.String("id", id))

	db, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if err = s.r.Delete(ctx, id); err != nil {
		err = fmt.Errorf("не удалось удалить базу данных: %s", err)
		zap.S().Error(err, zap.String("id", id))
		return err
	}

	db.Close()
	delete(s.dbList, db.ID)

	zap.S().Info("база данных успешно удалена", zap.String("id", id))
	return nil
}

func (s *service) TableList(ctx context.Context, id string) ([]*dbmodel.Table, error) {
	zap.S().Info("попытка получить таблицы базы данных", zap.String("id", id))

	db, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	var list []*dbmodel.Table
	if list, err = db.TableList(ctx); err != nil {
		err = fmt.Errorf("не удалось получить таблицы базы данных: %s", err)
		zap.S().Error(err, zap.String("id", id))
		return nil, err
	}

	zap.S().Info("таблицы базы данных успешно получены", zap.String("id", id))
	return list, nil
}

func (s *service) FunctionList(ctx context.Context, id string) ([]*dbmodel.Function, error) {
	zap.S().Info("попытка получить функции базы данных", zap.String("id", id))

	db, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	var list []*dbmodel.Function
	if list, err = db.FunctionList(ctx); err != nil {
		err = fmt.Errorf("не удалось получить функции базы данных: %s", err)
		zap.S().Error(err, zap.String("id", id))
		return nil, err
	}

	zap.S().Info("функции базы данных успешно получено", zap.String("id", id))
	return list, nil
}

func New(r DBRepo, tx database.TxManager) (*service, error) {
	s := &service{r: r, tx: tx, dbList: make(map[string]*dbmodel.DB)}

	list, err := r.GetList(context.Background())
	if err != nil {
		err = fmt.Errorf("не удалось получить базы данных из базы данных: %s", err)
		zap.S().Error(err)
		return nil, err
	}

	for _, d := range list {
		if err = d.Open(); err != nil {
			zap.S().Infof("не удалось подключиться к ранее добавленной базе данных: %s", err)
		}
		s.dbList[d.ID] = d
	}

	return s, nil
}
