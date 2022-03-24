// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent/kind"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent/predicate"
)

// KindQuery is the builder for querying Kind entities.
type KindQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.Kind
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the KindQuery builder.
func (kq *KindQuery) Where(ps ...predicate.Kind) *KindQuery {
	kq.predicates = append(kq.predicates, ps...)
	return kq
}

// Limit adds a limit step to the query.
func (kq *KindQuery) Limit(limit int) *KindQuery {
	kq.limit = &limit
	return kq
}

// Offset adds an offset step to the query.
func (kq *KindQuery) Offset(offset int) *KindQuery {
	kq.offset = &offset
	return kq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (kq *KindQuery) Unique(unique bool) *KindQuery {
	kq.unique = &unique
	return kq
}

// Order adds an order step to the query.
func (kq *KindQuery) Order(o ...OrderFunc) *KindQuery {
	kq.order = append(kq.order, o...)
	return kq
}

// First returns the first Kind entity from the query.
// Returns a *NotFoundError when no Kind was found.
func (kq *KindQuery) First(ctx context.Context) (*Kind, error) {
	nodes, err := kq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{kind.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (kq *KindQuery) FirstX(ctx context.Context) *Kind {
	node, err := kq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Kind ID from the query.
// Returns a *NotFoundError when no Kind ID was found.
func (kq *KindQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = kq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{kind.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (kq *KindQuery) FirstIDX(ctx context.Context) int {
	id, err := kq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Kind entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Kind entity is found.
// Returns a *NotFoundError when no Kind entities are found.
func (kq *KindQuery) Only(ctx context.Context) (*Kind, error) {
	nodes, err := kq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{kind.Label}
	default:
		return nil, &NotSingularError{kind.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (kq *KindQuery) OnlyX(ctx context.Context) *Kind {
	node, err := kq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Kind ID in the query.
// Returns a *NotSingularError when more than one Kind ID is found.
// Returns a *NotFoundError when no entities are found.
func (kq *KindQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = kq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{kind.Label}
	default:
		err = &NotSingularError{kind.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (kq *KindQuery) OnlyIDX(ctx context.Context) int {
	id, err := kq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Kinds.
func (kq *KindQuery) All(ctx context.Context) ([]*Kind, error) {
	if err := kq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return kq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (kq *KindQuery) AllX(ctx context.Context) []*Kind {
	nodes, err := kq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Kind IDs.
func (kq *KindQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := kq.Select(kind.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (kq *KindQuery) IDsX(ctx context.Context) []int {
	ids, err := kq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (kq *KindQuery) Count(ctx context.Context) (int, error) {
	if err := kq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return kq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (kq *KindQuery) CountX(ctx context.Context) int {
	count, err := kq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (kq *KindQuery) Exist(ctx context.Context) (bool, error) {
	if err := kq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return kq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (kq *KindQuery) ExistX(ctx context.Context) bool {
	exist, err := kq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the KindQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (kq *KindQuery) Clone() *KindQuery {
	if kq == nil {
		return nil
	}
	return &KindQuery{
		config:     kq.config,
		limit:      kq.limit,
		offset:     kq.offset,
		order:      append([]OrderFunc{}, kq.order...),
		predicates: append([]predicate.Kind{}, kq.predicates...),
		// clone intermediate query.
		sql:    kq.sql.Clone(),
		path:   kq.path,
		unique: kq.unique,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Kind.Query().
//		GroupBy(kind.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (kq *KindQuery) GroupBy(field string, fields ...string) *KindGroupBy {
	group := &KindGroupBy{config: kq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := kq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return kq.sqlQuery(ctx), nil
	}
	return group
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.Kind.Query().
//		Select(kind.FieldName).
//		Scan(ctx, &v)
//
func (kq *KindQuery) Select(fields ...string) *KindSelect {
	kq.fields = append(kq.fields, fields...)
	return &KindSelect{KindQuery: kq}
}

func (kq *KindQuery) prepareQuery(ctx context.Context) error {
	for _, f := range kq.fields {
		if !kind.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if kq.path != nil {
		prev, err := kq.path(ctx)
		if err != nil {
			return err
		}
		kq.sql = prev
	}
	return nil
}

func (kq *KindQuery) sqlAll(ctx context.Context) ([]*Kind, error) {
	var (
		nodes = []*Kind{}
		_spec = kq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &Kind{config: kq.config}
		nodes = append(nodes, node)
		return node.scanValues(columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		return node.assignValues(columns, values)
	}
	if err := sqlgraph.QueryNodes(ctx, kq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (kq *KindQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := kq.querySpec()
	_spec.Node.Columns = kq.fields
	if len(kq.fields) > 0 {
		_spec.Unique = kq.unique != nil && *kq.unique
	}
	return sqlgraph.CountNodes(ctx, kq.driver, _spec)
}

func (kq *KindQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := kq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (kq *KindQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   kind.Table,
			Columns: kind.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: kind.FieldID,
			},
		},
		From:   kq.sql,
		Unique: true,
	}
	if unique := kq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := kq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, kind.FieldID)
		for i := range fields {
			if fields[i] != kind.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := kq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := kq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := kq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := kq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (kq *KindQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(kq.driver.Dialect())
	t1 := builder.Table(kind.Table)
	columns := kq.fields
	if len(columns) == 0 {
		columns = kind.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if kq.sql != nil {
		selector = kq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if kq.unique != nil && *kq.unique {
		selector.Distinct()
	}
	for _, p := range kq.predicates {
		p(selector)
	}
	for _, p := range kq.order {
		p(selector)
	}
	if offset := kq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := kq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// KindGroupBy is the group-by builder for Kind entities.
type KindGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (kgb *KindGroupBy) Aggregate(fns ...AggregateFunc) *KindGroupBy {
	kgb.fns = append(kgb.fns, fns...)
	return kgb
}

// Scan applies the group-by query and scans the result into the given value.
func (kgb *KindGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := kgb.path(ctx)
	if err != nil {
		return err
	}
	kgb.sql = query
	return kgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (kgb *KindGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := kgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (kgb *KindGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(kgb.fields) > 1 {
		return nil, errors.New("ent: KindGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := kgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (kgb *KindGroupBy) StringsX(ctx context.Context) []string {
	v, err := kgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (kgb *KindGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = kgb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{kind.Label}
	default:
		err = fmt.Errorf("ent: KindGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (kgb *KindGroupBy) StringX(ctx context.Context) string {
	v, err := kgb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (kgb *KindGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(kgb.fields) > 1 {
		return nil, errors.New("ent: KindGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := kgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (kgb *KindGroupBy) IntsX(ctx context.Context) []int {
	v, err := kgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (kgb *KindGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = kgb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{kind.Label}
	default:
		err = fmt.Errorf("ent: KindGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (kgb *KindGroupBy) IntX(ctx context.Context) int {
	v, err := kgb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (kgb *KindGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(kgb.fields) > 1 {
		return nil, errors.New("ent: KindGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := kgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (kgb *KindGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := kgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (kgb *KindGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = kgb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{kind.Label}
	default:
		err = fmt.Errorf("ent: KindGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (kgb *KindGroupBy) Float64X(ctx context.Context) float64 {
	v, err := kgb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (kgb *KindGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(kgb.fields) > 1 {
		return nil, errors.New("ent: KindGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := kgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (kgb *KindGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := kgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (kgb *KindGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = kgb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{kind.Label}
	default:
		err = fmt.Errorf("ent: KindGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (kgb *KindGroupBy) BoolX(ctx context.Context) bool {
	v, err := kgb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (kgb *KindGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range kgb.fields {
		if !kind.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := kgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := kgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (kgb *KindGroupBy) sqlQuery() *sql.Selector {
	selector := kgb.sql.Select()
	aggregation := make([]string, 0, len(kgb.fns))
	for _, fn := range kgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(kgb.fields)+len(kgb.fns))
		for _, f := range kgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(kgb.fields...)...)
}

// KindSelect is the builder for selecting fields of Kind entities.
type KindSelect struct {
	*KindQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (ks *KindSelect) Scan(ctx context.Context, v interface{}) error {
	if err := ks.prepareQuery(ctx); err != nil {
		return err
	}
	ks.sql = ks.KindQuery.sqlQuery(ctx)
	return ks.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ks *KindSelect) ScanX(ctx context.Context, v interface{}) {
	if err := ks.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (ks *KindSelect) Strings(ctx context.Context) ([]string, error) {
	if len(ks.fields) > 1 {
		return nil, errors.New("ent: KindSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := ks.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ks *KindSelect) StringsX(ctx context.Context) []string {
	v, err := ks.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (ks *KindSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = ks.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{kind.Label}
	default:
		err = fmt.Errorf("ent: KindSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (ks *KindSelect) StringX(ctx context.Context) string {
	v, err := ks.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (ks *KindSelect) Ints(ctx context.Context) ([]int, error) {
	if len(ks.fields) > 1 {
		return nil, errors.New("ent: KindSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := ks.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ks *KindSelect) IntsX(ctx context.Context) []int {
	v, err := ks.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (ks *KindSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = ks.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{kind.Label}
	default:
		err = fmt.Errorf("ent: KindSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (ks *KindSelect) IntX(ctx context.Context) int {
	v, err := ks.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (ks *KindSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(ks.fields) > 1 {
		return nil, errors.New("ent: KindSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := ks.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ks *KindSelect) Float64sX(ctx context.Context) []float64 {
	v, err := ks.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (ks *KindSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = ks.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{kind.Label}
	default:
		err = fmt.Errorf("ent: KindSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (ks *KindSelect) Float64X(ctx context.Context) float64 {
	v, err := ks.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (ks *KindSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(ks.fields) > 1 {
		return nil, errors.New("ent: KindSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := ks.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ks *KindSelect) BoolsX(ctx context.Context) []bool {
	v, err := ks.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (ks *KindSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = ks.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{kind.Label}
	default:
		err = fmt.Errorf("ent: KindSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (ks *KindSelect) BoolX(ctx context.Context) bool {
	v, err := ks.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ks *KindSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ks.sql.Query()
	if err := ks.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}