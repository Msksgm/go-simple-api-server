package postgres

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

func formatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	} else if offset > 0 {
		return fmt.Sprintf("OFFSET %d", offset)
	}
	return ""
}

func formatWhereClause(where []string) string {
	if len(where) == 0 {
		return ""
	}
	return " WHERE " + strings.Join(where, " AND ")
}

func findMany(ctx context.Context, tx *sqlx.Tx, ss interface{}, query string, args ...interface{}) error {
	rows, err := tx.QueryxContext(ctx, query, args...)
	if err != nil {
		return err
	}

	defer rows.Close()

	sPtrVal, err := asSlicePtrValue(ss)
	if err != nil {
		return err
	}

	sVal := sPtrVal.Elem()
	newSlice := reflect.MakeSlice(sVal.Type(), 0, 0)
	elemType := sliceElemType(sVal)

	for rows.Next() {
		newVal := reflect.New(elemType)
		if err := rows.StructScan(newVal.Interface()); err != nil {
			return nil
		}
		newSlice = reflect.Append(newSlice, newVal)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	sPtrVal.Elem().Set(newSlice)

	return nil
}

func sliceElemType(v reflect.Value) reflect.Type {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	vv := v.Type().Elem()

	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	return vv
}

func isSlicePtr(v interface{}) bool {
	typ := reflect.TypeOf(v)

	return typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Slice
}

func asSlicePtrValue(v interface{}) (reflect.Value, error) {
	if !isSlicePtr(v) {
		return reflect.Value{}, errors.New("expecting a pointer to slice")
	}
	return reflect.ValueOf(v), nil
}
