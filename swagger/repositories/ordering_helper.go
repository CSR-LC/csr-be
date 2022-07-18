package repositories

import (
	"errors"
	"fmt"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
)

const (
	ascOrder  string = "asc"
	descOrder string = "desc"
)

func getOrderFunc(orderBy, orderColumn string) (ent.OrderFunc, error) {
	var orderFunc ent.OrderFunc
	var err error
	switch orderBy {
	case ascOrder:
		orderFunc = ent.Asc(orderColumn)
	case descOrder:
		orderFunc = ent.Desc(orderColumn)
	default:
		err = errors.New(fmt.Sprintf("wrong value for orderBy: %s", orderBy))
	}
	return orderFunc, err
}

func checkOrderColumn(orderColumn string, fields []string) bool {
	for _, f := range fields {
		if orderColumn == f {
			return true
		}
	}
	return false
}
