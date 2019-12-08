package manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/shelton-hu/util/jsonutil"
	"github.com/shelton-hu/util/stringutil"

	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/logger"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/pi"
)

var (
	ErrUnknownPrismaClient = errors.New("Unknown prisma client")
)

func BuildQueryOrderDir(req Request, defaultOrderDir string) string {
	orderDir := defaultOrderDir
	if rws, ok := req.(RequestWithSort); ok {
		r := rws.GetReverse()
		s := rws.GetSortKey()
		if s == nil {
			return orderDir
		}
		reverseKey := "DESC"
		reverse := r.GetValue()
		if reverse {
			reverseKey = "ASC"
		}
		sortKey := s.GetValue()
		orderDir = stringutil.StringJoin(sortKey, "_", reverseKey)
	}
	return orderDir
}

func NativeGQL(ctx context.Context, endpoint string, receipt interface{}, query string, params ...interface{}) error {
	var result map[string]interface{}
	var err error
	var queryString = fmt.Sprintf(query, params...)
	_ = queryString

	result, err = pi.Global().MysqlPrisma(ctx).GraphQL(ctx, queryString, make(map[string]interface{}))
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}

	b, err := jsonutil.Encode(result)
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}
	return jsonutil.Decode(b, &receipt)
}
