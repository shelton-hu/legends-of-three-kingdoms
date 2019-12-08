package iam

import (
	"context"
)

func (s *Server) Checker(ctx context.Context, req interface{}) error {
	return nil
}

func (s *Server) Builder(ctx context.Context, req interface{}) interface{} {
	return req
}
