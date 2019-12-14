package iam

import (
	"context"
	"errors"
	"time"

	"github.com/shelton-hu/util/hashutil"
	"github.com/shelton-hu/util/pbutil"
	"github.com/shelton-hu/util/pointerutil"
	"github.com/shelton-hu/util/timeutil"

	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/gerr"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/logger"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/pb"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/pi"
	prisma "github.com/shelton-hu/legends-of-three-kingdoms/pkg/prisma/mysql-prisma-client"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/util/senderutil"
)

var (
	ErrAlreadyInGame = errors.New("您还在游戏中")
	ErrorPasseord    = errors.New("密码错误")
)

func (s *Server) SignInOrSignUp(ctx context.Context, req *pb.SignInOrSignUpRequest) (*pb.SignInOrSignUpResponse, error) {
	// 1.获取入参
	nickName := req.GetNickName().GetValue()
	password := req.GetPassword().GetValue()

	// 2.查询用户是否存在
	isNewUser := false
	user, err := pi.Global().MysqlPrisma(ctx).User(prisma.UserWhereUniqueInput{
		NickName: &nickName,
	}).Exec(ctx)
	if err != nil {
		if err.Error() == prisma.ErrNoResult.Error() {
			isNewUser = true
			logger.Debug(ctx, "create a new user, %+v", isNewUser)
		} else {
			logger.Error(ctx, err.Error())
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorSignInOrSingUpFailed, nickName)
		}
	} else {
		if user.IsInGame {
			logger.Error(ctx, ErrAlreadyInGame.Error())
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, ErrAlreadyInGame, gerr.ErrorSignInOrSingUpFailed, nickName)
		}
	}

	// 3.登录验证或注册
	if isNewUser {
		// 3.1注册
		user, err = pi.Global().MysqlPrisma(ctx).CreateUser(prisma.UserCreateInput{
			NickName: nickName,
			Password: pointerutil.GetStringPointer(hashutil.GetStringMd5(password)),
		}).Exec(ctx)
		if err != nil {
			logger.Error(ctx, err.Error())
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorSignInOrSingUpFailed, nickName)
		}
	} else {
		// 3.2登录验证
		if pointerutil.GetString(user.Password) != hashutil.GetStringMd5(password) {
			logger.Error(ctx, ErrorPasseord.Error())
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, ErrorPasseord, gerr.ErrorSignInOrSingUpFailed, nickName)
		}
	}

	// 4.更新"token"及"loginAt"
	accessToken, err := senderutil.Generate(
		pi.Global().Cfg(ctx).IAM.SecretKey, pi.Global().Cfg(ctx).IAM.ExpireTime, user.ID,
	)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorSignInOrSingUpFailed, nickName)
	}
	_, err = pi.Global().MysqlPrisma(ctx).UpdateUser(prisma.UserUpdateParams{
		Data: prisma.UserUpdateInput{
			LoginAt: pointerutil.GetStringPointer(timeutil.Time2String(time.Now(), time.RFC3339)),
			Token:   pointerutil.GetStringPointer(hashutil.GetStringMd5(accessToken)),
		},
		Where: prisma.UserWhereUniqueInput{
			NickName: &nickName,
		},
	}).Exec(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorSignInOrSingUpFailed, nickName)
	}

	return &pb.SignInOrSignUpResponse{
		AccessToken: pbutil.ToProtoString(accessToken),
	}, nil
}

func (s *Server) SignOut(ctx context.Context, req *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	return nil, nil
}
