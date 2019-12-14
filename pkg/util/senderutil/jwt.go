package senderutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/logger"
)

var ErrExpired = fmt.Errorf("access token expired")

func trimKey(k string) []byte {
	return []byte(strings.TrimSpace(k))
}

func Validate(k, s string) (*Sender, error) {
	ctx := context.Background()
	tok, err := jwt.ParseSigned(s)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, err
	}
	c := &jwt.Claims{}
	sender := &Sender{}
	err = tok.Claims(trimKey(k), c, sender)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, err
	}
	if c.Expiry.Time().Unix() < time.Now().Unix() {
		logger.Error(ctx, ErrExpired.Error())
		return nil, ErrExpired
	}
	sender.UserId = c.Subject
	return sender, nil
}

func Generate(k string, expire time.Duration, userId string) (string, error) {
	ctx := context.Background()
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS512, Key: trimKey(k)}, nil)
	if err != nil {
		logger.Error(ctx, ErrExpired.Error())
		return "", err
	}
	sender := &Sender{}
	now := time.Now()
	c := &jwt.Claims{
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(expire)),
		Subject:  userId,
	}
	return jwt.Signed(signer).Claims(sender).Claims(c).CompactSerialize()
}
