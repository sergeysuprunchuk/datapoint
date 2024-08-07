package querymodel

import (
	"datapoint/internal/model/dbmodel"
	sq "github.com/Masterminds/squirrel"
	"reflect"
	"testing"
)

type test struct {
	query         Query
	expectedQuery string
	expectedArgs  []any
}

var b = sq.
	StatementBuilder.
	PlaceholderFormat(sq.Question)

var table = &Table{TableKey: TableKey{Name: "example"}}

func TestBuildSelect(t *testing.T) {
	tests := [...]test{
		{
			query: Query{
				Type: Select,
				Table: &Table{
					TableKey: table.TableKey,
					Next: []*Table{
						{
							TableKey: TableKey{Name: table.Name, Increment: 1},
							Rule: &Rule{
								Type: Join,
								Conditions: []*Condition{{
									Columns: [2]*Column{
										{
											Column:   dbmodel.Column{Name: "id"},
											TableKey: table.TableKey,
										},
										{
											Column:   dbmodel.Column{Name: "id"},
											TableKey: TableKey{Name: table.TableKey.Name, Increment: 1},
										},
									},
									Operator: Equal,
								}},
							},
						},
						{
							TableKey: TableKey{Name: table.Name, Increment: 2},
							Next: []*Table{
								{
									TableKey: TableKey{Name: table.Name, Increment: 3},
									Rule: &Rule{
										Type: Right,
										Conditions: []*Condition{{
											Columns: [2]*Column{
												{
													Column:   dbmodel.Column{Name: "id"},
													TableKey: TableKey{Name: table.TableKey.Name, Increment: 2},
												},
												{
													Column:   dbmodel.Column{Name: "id"},
													TableKey: TableKey{Name: table.TableKey.Name, Increment: 3},
												},
											},
											Operator: Equal,
										}},
									},
								},
							},
							Rule: &Rule{
								Type: Left,
								Conditions: []*Condition{{
									Columns: [2]*Column{
										{
											Column:   dbmodel.Column{Name: "id"},
											TableKey: table.TableKey,
										},
										{
											Column:   dbmodel.Column{Name: "id"},
											TableKey: TableKey{Name: table.TableKey.Name, Increment: 2},
										},
									},
									Operator: Equal,
								}},
							},
						},
					},
				},
				Columns: []*Column{
					{
						TableKey: table.TableKey,
						Column:   dbmodel.Column{Name: "id"},
					},
					{
						TableKey: table.TableKey,
						Column:   dbmodel.Column{Name: "name"},
					},
					{
						TableKey: table.TableKey,
						Column:   dbmodel.Column{Name: "age"},
						Function: "sum",
					},
					{
						TableKey: TableKey{Name: table.TableKey.Name, Increment: 3},
						Column:   dbmodel.Column{Name: "age"},
						Function: "avg",
					},
					{
						TableKey: TableKey{Name: table.TableKey.Name, Increment: 1},
						Column:   dbmodel.Column{Name: "id"},
					},
				},
				OrderBy: []*Column{
					{
						Column:   dbmodel.Column{Name: "name"},
						TableKey: TableKey{Name: table.TableKey.Name, Increment: 1},
						Desc:     true,
					},
				},
				b: b,
			},
			expectedQuery: "SELECT \"example\".\"id\" \"example.id\", " +
				"\"example\".\"name\" \"example.name\", " +
				"sum(\"example\".\"age\") \"sum(example.age)\", " +
				"avg(\"example3\".\"age\") \"avg(example3.age)\", " +
				"\"example1\".\"id\" \"example1.id\" " +
				"FROM \"example\" \"example\" " +
				"JOIN \"example\" \"example1\" ON \"example\".\"id\" = \"example1\".\"id\" " +
				"LEFT JOIN \"example\" \"example2\" ON \"example\".\"id\" = \"example2\".\"id\" " +
				"RIGHT JOIN \"example\" \"example3\" ON \"example2\".\"id\" = \"example3\".\"id\" " +
				"GROUP BY \"example\".\"id\", \"example\".\"name\", \"example1\".\"id\" " +
				"ORDER BY \"example1\".\"name\" DESC",
		},
	}

	for _, test := range tests {
		query, args, err := test.query.buildSelect().ToSql()
		if err != nil {
			t.Errorf("произошла ошибка при построении запроса: %s", err)
		}

		if query != test.expectedQuery || !reflect.DeepEqual(args, test.expectedArgs) {
			t.Errorf(`query --> ожидалось: %s, получено: %s;
args --> ожидалось: %v, получено: %v`, test.expectedQuery, query, test.expectedArgs, args)
		}
	}
}

func TestBuildInsert(t *testing.T) {
	tests := [...]test{
		{
			query: Query{
				Type:  Insert,
				Table: table,
				Columns: []*Column{
					{
						Column: dbmodel.Column{Name: "id"},
						Value:  "slvag",
					},
					{
						Column: dbmodel.Column{Name: "name"},
						Value:  "qtbbt",
					},
					{
						Column: dbmodel.Column{Name: "age"},
						Value:  70,
					},
				},
				b: b,
			},
			expectedQuery: `INSERT INTO "example" ("id","name","age") VALUES (?,?,?)`,
			expectedArgs:  []any{"slvag", "qtbbt", 70},
		},
	}

	for _, test := range tests {
		query, args, err := test.query.buildInsert().ToSql()
		if err != nil {
			t.Errorf("произошла ошибка при построении запроса: %s", err)
		}

		if query != test.expectedQuery || !reflect.DeepEqual(args, test.expectedArgs) {
			t.Errorf(`query --> ожидалось: %s, получено: %s;
args --> ожидалось: %v, получено: %v`, test.expectedQuery, query, test.expectedArgs, args)
		}
	}
}

func TestBuildUpdate(t *testing.T) {
	tests := [...]test{
		{
			query: Query{
				Type:  Update,
				Table: table,
				Columns: []*Column{
					{
						Column: dbmodel.Column{Name: "name"},
						Value:  "qtbbt",
					},
					{
						Column: dbmodel.Column{Name: "age"},
						Value:  70,
					},
				},
				Where: []*Column{
					{
						Column: dbmodel.Column{Name: "id"},
						Value:  "slvag",
					},
					{
						Column: dbmodel.Column{Name: "name"},
						Value:  "qtbbt",
					},
				},
				b: b,
			},
			expectedQuery: `UPDATE "example" SET "name" = ?, "age" = ? WHERE "id" = ? AND "name" = ?`,
			expectedArgs:  []any{"qtbbt", 70, "slvag", "qtbbt"},
		},
	}

	for _, test := range tests {
		query, args, err := test.query.buildUpdate().ToSql()
		if err != nil {
			t.Errorf("произошла ошибка при построении запроса: %s", err)
		}

		if query != test.expectedQuery || !reflect.DeepEqual(args, test.expectedArgs) {
			t.Errorf(`query --> ожидалось: %s, получено: %s;
args --> ожидалось: %v, получено: %v`, test.expectedQuery, query, test.expectedArgs, args)
		}
	}
}

func TestBuildDelete(t *testing.T) {
	tests := [...]test{
		{
			query: Query{
				Type:  Delete,
				Table: table,
				Where: []*Column{
					{
						Column: dbmodel.Column{Name: "id"},
						Value:  "slvag",
					},
				},
				b: b,
			},
			expectedQuery: `DELETE FROM "example" WHERE "id" = ?`,
			expectedArgs:  []any{"slvag"},
		},
		{
			query: Query{
				Type:  Delete,
				Table: table,
				Where: []*Column{
					{
						Column: dbmodel.Column{Name: "id"},
						Value:  "slvag",
					},
					{
						Column: dbmodel.Column{Name: "name"},
						Value:  "qtbbt",
					},
				},
				b: b,
			},
			expectedQuery: `DELETE FROM "example" WHERE "id" = ? AND "name" = ?`,
			expectedArgs:  []any{"slvag", "qtbbt"},
		},
	}

	for _, test := range tests {
		query, args, err := test.query.buildDelete().ToSql()
		if err != nil {
			t.Errorf("произошла ошибка при построении запроса: %s", err)
		}

		if query != test.expectedQuery || !reflect.DeepEqual(args, test.expectedArgs) {
			t.Errorf(`query --> ожидалось: %s, получено: %s;
args --> ожидалось: %v, получено: %v`, test.expectedQuery, query, test.expectedArgs, args)
		}
	}
}
