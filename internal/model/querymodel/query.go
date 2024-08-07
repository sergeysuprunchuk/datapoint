package querymodel

import (
	"context"
	"database/sql"
	"datapoint/internal/model/dbmodel"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"strconv"
	"strings"
)

const (
	Select = "select"
	Insert = "insert"
	Update = "update"
	Delete = "delete"
)

type Runner interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

type Query struct {
	Type    string
	Table   *Table
	Columns []*Column
	OrderBy []*Column
	Where   []*Column
	Limit   uint64
	Offset  uint64
	b       sq.StatementBuilderType //с PlaceholderFormat, но без RunWith
}

func (q Query) Execute(ctx context.Context, runner Runner) (QueryResult, error) {
	switch q.Type {
	case Select:
		return q.executeSelect(ctx, runner)
	case Insert:
		return q.executeInsert(ctx, runner)
	case Update:
		return q.executeUpdate(ctx, runner)
	case Delete:
		return q.executeDelete(ctx, runner)
	default:
		return QueryResult{}, fmt.Errorf("%s", q.Type)
	}
}

func (q Query) buildSelect() sq.SelectBuilder {
	b := q.b.
		Select().
		From(fmt.Sprintf(`"%s" "%s"`, q.Table.Name, q.Table))

	next := q.Table.Next
	for i := 0; i < len(next); i++ {
		rule := fmt.Sprintf(`"%s" "%s" ON %s`, next[i].Name, next[i].TableKey, next[i].Rule.String())
		switch next[i].Rule.Type {
		case Join:
			b = b.Join(rule)
		case Left:
			b = b.LeftJoin(rule)
		case Right:
			b = b.RightJoin(rule)
		}
		if next[i].Next != nil {
			next = append(next, next[i].Next...)
		}
	}

	var (
		groupBy     []string
		hasFunction bool
	)

	for _, c := range q.Columns {
		b = b.Columns(c.StringWTWA())
		if len(c.Function) != 0 {
			hasFunction = true
			continue
		}
		groupBy = append(groupBy, c.StringWT())
	}

	for _, c := range q.OrderBy {
		var order string
		if c.Desc {
			order = " DESC"
		}
		b = b.OrderBy(c.StringWT() + order)
	}

	var where sq.Eq
	for _, c := range q.Where {
		if c.Value == nil {
			continue
		}

		if where == nil {
			where = make(sq.Eq, len(q.Where))
		}

		where[c.StringWT()] = c.Value
	}
	if where != nil {
		b = b.Where(where)
	}

	if hasFunction {
		b = b.GroupBy(groupBy...)
	}

	return b
}

func (q Query) executeSelect(ctx context.Context, runner Runner) (QueryResult, error) {
	query, args, err := q.buildSelect().ToSql()
	if err != nil {
		return QueryResult{}, err
	}

	var rows *sql.Rows
	if rows, err = runner.QueryContext(ctx, query, args...); err != nil {
		return QueryResult{}, err
	}
	defer func() { _ = rows.Close() }()

	var (
		data    []map[string]any
		columns []string
	)

	if columns, err = rows.Columns(); err != nil {
		return QueryResult{}, err
	}

	for rows.Next() {
		i, dest := make(map[string]any), make([]any, 0)

		for _, c := range columns {
			i[c] = new(any)
			dest = append(dest, i[c])
		}

		if err = rows.Scan(dest...); err != nil {
			return QueryResult{}, err
		}

		data = append(data, i)
	}

	return QueryResult{Data: data}, nil
}

func (q Query) buildInsert() sq.InsertBuilder {
	b := q.b.Insert(strconv.Quote(q.Table.Name))

	values := make([]any, 0, len(q.Columns))

	for _, c := range q.Columns {
		b = b.Columns(strconv.Quote(c.Name))
		values = append(values, c.Value)
	}

	return b.Values(values...)
}

func (q Query) executeInsert(ctx context.Context, runner Runner) (QueryResult, error) {
	query, args, err := q.buildInsert().ToSql()
	if err != nil {
		return QueryResult{}, err
	}

	_, err = runner.ExecContext(ctx, query, args...)

	return QueryResult{}, err
}

func (q Query) buildUpdate() sq.UpdateBuilder {
	b := q.b.Update(strconv.Quote(q.Table.Name))

	for _, c := range q.Columns {
		b = b.Set(strconv.Quote(c.Name), c.Value)
	}

	where := make(sq.Eq, len(q.Where))

	for _, c := range q.Where {
		where[strconv.Quote(c.Name)] = c.Value
	}

	return b.Where(where)
}

func (q Query) executeUpdate(ctx context.Context, runner Runner) (QueryResult, error) {
	query, args, err := q.buildUpdate().ToSql()
	if err != nil {
		return QueryResult{}, err
	}

	_, err = runner.ExecContext(ctx, query, args...)

	return QueryResult{}, err
}

func (q Query) buildDelete() sq.DeleteBuilder {
	b := q.b.Delete(strconv.Quote(q.Table.Name))

	where := make(sq.Eq, len(q.Where))

	for _, c := range q.Where {
		where[strconv.Quote(c.Name)] = c.Value
	}

	return b.Where(where)
}

func (q Query) executeDelete(ctx context.Context, runner Runner) (QueryResult, error) {
	query, args, err := q.buildDelete().ToSql()
	if err != nil {
		return QueryResult{}, err
	}

	_, err = runner.ExecContext(ctx, query, args...)

	return QueryResult{}, err
}

type TableKey struct {
	Name      string
	Increment uint8
}

func (k TableKey) String() string {
	if k.Increment != 0 {
		return fmt.Sprintf("%s%d", k.Name, k.Increment)
	}
	return k.Name
}

/*
суффиксы функций:
	WT - с таблицей
	WA - с псевдонимом
*/

type Column struct {
	dbmodel.Column
	TableKey TableKey
	Function string
	Desc     bool
	Value    any
}

func (c Column) String() string {
	if len(c.Function) != 0 {
		return fmt.Sprintf(`%s("%s")`, c.Function, c.Name)
	}
	return strconv.Quote(c.Name)
}

func (c Column) StringWT() string {
	if len(c.Function) != 0 {
		return fmt.Sprintf(`%s("%s"."%s")`, c.Function, c.TableKey, c.Name)
	}
	return fmt.Sprintf(`"%s"."%s"`, c.TableKey, c.Name)
}

func (c Column) StringWTWA() string {
	if len(c.Function) != 0 {
		return fmt.Sprintf(`%s "%s(%s.%s)"`, c.StringWT(), c.Function, c.TableKey, c.Name)
	}
	return fmt.Sprintf(`%s "%s.%s"`, c.StringWT(), c.TableKey, c.Name)
}

type Table struct {
	TableKey
	Next []*Table
	Rule *Rule
}

const (
	Left  = "left"
	Right = "right"
	Join  = "join"
)

type Rule struct {
	Type       string
	Conditions []*Condition
}

func (r Rule) String() string {
	and := make([]string, 0, len(r.Conditions))

	for _, c := range r.Conditions {
		and = append(and, fmt.Sprintf("%s %s %s", c.Columns[0].StringWT(), c.Operator, c.Columns[1].StringWT()))
	}

	return strings.Join(and, " AND ")
}

const (
	Equal    = "="
	NotEqual = "!="
)

type Condition struct {
	Columns  [2]*Column
	Operator string
}

type QueryResult struct{ Data []map[string]any }
