package senderutil

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/metadata"

	"github.com/shelton-hu/util/jsonutil"

	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/util/ctxutil"
)

const (
	SenderKey = "sender"
	TokenType = "Bearer"
)

type Sender struct {
	UserId   string `json:"user_id,omitempty"`
	UserType string `json:"user_type,omitempty"`
	Role     string `json:"role,omitempty"`
}

func (s *Sender) ToJson() string {
	return jsonutil.ToString(s)
}

func GetSenderFromContext(ctx context.Context) *Sender {
	values := ctxutil.GetValueFromContext(ctx, SenderKey)
	if len(values) == 0 || len(values[0]) == 0 {
		return nil
	}
	sender := Sender{}
	err := json.Unmarshal([]byte(values[0]), &sender)
	if err != nil {
		panic(err)
	}
	return &sender
}

func ContextWithSender(ctx context.Context, user *Sender) context.Context {
	if user == nil {
		return ctx
	}
	ctx = context.WithValue(ctx, SenderKey, []string{user.ToJson()})
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	md[SenderKey] = []string{user.ToJson()}
	return metadata.NewOutgoingContext(ctx, md)
}
