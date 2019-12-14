package room

import (
	"context"
	"errors"

	"github.com/shelton-hu/util/pbutil"
	"github.com/shelton-hu/util/pointerutil"

	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/gerr"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/logger"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/pb"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/pi"
	prisma "github.com/shelton-hu/legends-of-three-kingdoms/pkg/prisma/mysql-prisma-client"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/util/senderutil"
)

var (
	ErrAlreadyInGame = errors.New("您还在游戏中")
)

func (s *Server) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {
	// 1.获取入参
	roomNickName := req.GetRoomNickName().GetValue()
	userId := senderutil.GetSenderFromContext(ctx).UserId

	// 2.校验玩家状态
	if err := s.validatePlayerStatus(ctx, userId); err != nil {
		logger.Error(ctx, err.Error())
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorComeIntoRoomFailed)
	}

	// 3.操作数据库--创建"Room"
	room, err := pi.Global().MysqlPrisma(ctx).CreateRoom(prisma.RoomCreateInput{
		RoomNickName: &roomNickName,
		Players: &prisma.UserCreateManyWithoutRoomInput{
			Connect: []prisma.UserWhereUniqueInput{
				prisma.UserWhereUniqueInput{
					ID: &userId,
				},
			},
		},
	}).Exec(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateRoomFailed)
	}

	// 4.操作数据库--更新"User"
	_, err = pi.Global().MysqlPrisma(ctx).UpdateUser(prisma.UserUpdateParams{
		Data: prisma.UserUpdateInput{
			IsInGame: pointerutil.GetBoolPointer(true),
		},
		Where: prisma.UserWhereUniqueInput{
			ID: &userId,
		},
	}).Exec(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateRoomFailed)
	}

	return &pb.CreateRoomResponse{
		RoomId: pbutil.ToProtoString(room.ID),
	}, nil
}

func (s *Server) ComeIntoRoom(ctx context.Context, req *pb.ComeIntoRoomRequest) (*pb.ComeIntoRoomResponse, error) {
	// 1.获取入参
	roomId := req.GetRoomId()
	userId := senderutil.GetSenderFromContext(ctx).UserId

	// 2.校验玩家状态
	if err := s.validatePlayerStatus(ctx, userId); err != nil {
		logger.Error(ctx, err.Error())
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorComeIntoRoomFailed)
	}

	// 3.操作数据库--更新"Room"
	_, err := pi.Global().MysqlPrisma(ctx).UpdateRoom(prisma.RoomUpdateParams{
		Data: prisma.RoomUpdateInput{
			Players: &prisma.UserUpdateManyWithoutRoomInput{
				Connect: []prisma.UserWhereUniqueInput{
					prisma.UserWhereUniqueInput{
						ID: &userId,
					},
				},
			},
		},
		Where: prisma.RoomWhereUniqueInput{
			ID: &roomId,
		},
	}).Exec(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorComeIntoRoomFailed)
	}

	// 4.操作数据库--更新"User"
	_, err = pi.Global().MysqlPrisma(ctx).UpdateUser(prisma.UserUpdateParams{
		Data: prisma.UserUpdateInput{
			IsInGame: pointerutil.GetBoolPointer(true),
		},
		Where: prisma.UserWhereUniqueInput{
			ID: &userId,
		},
	}).Exec(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorComeIntoRoomFailed)
	}

	return &pb.ComeIntoRoomResponse{
		RoomId: pbutil.ToProtoString(roomId),
	}, nil
}

func (s *Server) DescribeRooms(ctx context.Context, req *pb.DescribeRoomsRequest) (*pb.DescribeRoomsResponse, error) {
	// 1.操作数据库
	results, err := pi.Global().MysqlPrisma(ctx).Rooms(nil).Exec(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeRoomsFailed)
	}

	// 2.组合出参
	rooms := make([]*pb.Room, len(results))
	for i, result := range results {
		rooms[i] = &pb.Room{
			RoomId:       pbutil.ToProtoString(result.ID),
			RoomNickName: pbutil.ToProtoString(pointerutil.GetString(result.RoomNickName)),
		}
	}

	return &pb.DescribeRoomsResponse{
		Rooms: rooms,
	}, nil
}

func (s *Server) validatePlayerStatus(ctx context.Context, userId string) error {
	user, err := pi.Global().MysqlPrisma(ctx).User(prisma.UserWhereUniqueInput{
		ID: &userId,
	}).Exec(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}
	if user.IsInGame {
		logger.Error(ctx, ErrAlreadyInGame.Error())
		return ErrAlreadyInGame
	}
	return nil
}
