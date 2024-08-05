package converter

import (
	"datapoint/internal/controller/http/model"
	"datapoint/internal/model/dbmodel"
	"datapoint/pkg/slices"
)

func ToDBInfo(i dbmodel.Info) model.DBInfo {
	return model.DBInfo{
		Name:     i.Name,
		Host:     i.Config.Host,
		Port:     int(i.Config.Port),
		DBUser:   i.Config.User,
		Password: i.Config.Password,
		DBName:   i.Config.Name,
		Driver:   i.Config.Driver,
	}
}

func ToDB(d *dbmodel.DB) model.DB {
	return model.DB{
		ID:     d.ID,
		DBInfo: ToDBInfo(d.Info),
	}
}

func ToDBList(list []*dbmodel.DB) []model.DB {
	return slices.Map(list, ToDB)
}

func FromDBInfo(i model.DBInfo) dbmodel.Info {
	return dbmodel.Info{
		Name: i.Name,
		Config: dbmodel.Config{
			Host:     i.Host,
			Port:     uint16(i.Port),
			User:     i.DBUser,
			Password: i.Password,
			Name:     i.DBName,
			Driver:   i.Driver,
		},
	}
}

func ToDBfk(fk *dbmodel.FK) *model.DBfk {
	if fk == nil {
		return nil
	}

	return &model.DBfk{
		TableName:  fk.TableName,
		ColumnName: fk.ColumnName,
	}
}

func ToDBColumn(c *dbmodel.Column) model.DBColumn {
	return model.DBColumn{
		Name:       c.Name,
		Type:       c.Type,
		IsRequired: c.IsRequired,
		IsPK:       c.IsPK,
		FK:         ToDBfk(c.FK),
	}
}

func ToDBColumnList(list []*dbmodel.Column) []model.DBColumn {
	return slices.Map(list, ToDBColumn)
}

func ToDBTable(t *dbmodel.Table) model.DBTable {
	return model.DBTable{
		Name:       t.Name,
		ColumnList: ToDBColumnList(t.ColumnList),
	}
}

func ToDBTableList(list []*dbmodel.Table) []model.DBTable {
	return slices.Map(list, ToDBTable)
}

func ToDBFunction(f *dbmodel.Function) model.DBFunction {
	return model.DBFunction{
		Name:     f.Name,
		TypeList: f.TypeList,
	}
}

func ToDBFunctionList(list []*dbmodel.Function) []model.DBFunction {
	return slices.Map(list, ToDBFunction)
}
