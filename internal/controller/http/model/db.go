package model

type DB struct {
	ID string `json:"id"`
	DBInfo
}

type DBInfo struct {
	Name     string `json:"name" validate:"required"`
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"min=1025,max=65535"`
	DBUser   string `json:"dbUser" validate:"required"`
	Password string `json:"password"`
	DBName   string `json:"dbName" validate:"required"`
	Driver   string `json:"driver" validate:"oneof=PostgreSQL"`
}

type DBfk struct {
	TableName  string `json:"tableName"`
	ColumnName string `json:"columnName"`
}

type DBColumn struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	IsRequired bool   `json:"isRequired"`
	IsPK       bool   `json:"isPK"`
	FK         *DBfk  `json:"fk,omitempty"`
}

type DBTable struct {
	Name       string     `json:"name"`
	ColumnList []DBColumn `json:"columnList"`
}

type DBFunction struct {
	Name     string   `json:"name"`
	TypeList []string `json:"typeList"`
}
