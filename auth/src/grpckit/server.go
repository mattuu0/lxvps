package grpckit

import (
	"auth/models"
	"context"
	"errors"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func Init() {
	log.Print("main start")

	// 9000番ポートでクライアントからのリクエストを受け付けるようにする
	listen, err := net.Listen("tcp", os.Getenv("GRPC_ADDR"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	RegisterAuthBaseServiceServer(grpcServer, &GrpcServer{})

	// 以下でリッスンし続ける
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	log.Print("main end")
}

type GrpcServer struct {

}

// GetLabel implements AuthBaseServiceServer.
func (grpcs *GrpcServer) GetLabel(ctx context.Context,req *GetLabelRequest) (*Label, error) {
	// ラベルを取得する関数
	// モデルから取得する
	getData,err := models.GetLabel(req.LabelName)

	// エラー処理
	if err != nil {
		return nil, err
	}

	// データを返す
	return &Label{
		Name:          getData.Name,
		Color:         getData.Color,
	}, nil
}

// GetUser implements AuthBaseServiceServer.
func (grpcs *GrpcServer) GetUser(ctx context.Context,req *GetUserRequest) (*User, error) {
	// ユーザーを取得する関数
	// モデルから取得する
	getData,rerr := models.GetUser(req.UserID)

	// エラー処理
	if rerr.Error != nil {
		return nil, rerr.Error
	}

	// データを返す
	return &User{
		UserID:        getData.UserID,
		Name:          getData.Name,
		Email:         getData.Email,
		Labels:        ModelLabelsToLabels(getData.Labels),
	}, nil
}

// SearchUser implements AuthBaseServiceServer.
func (grpcs *GrpcServer) SearchUser(ctx context.Context,req *SearchRequest) (*SearchResult, error) {
	// 名前を検索する関数

	// もし名前が指定されていたら
	if req.Name != "" {
		// 名前を検索する
		users, err := models.SearchUserByName(req.Name)
		if err != nil {
			return nil, err
		}

		return &SearchResult{
			Users: ModelUsersToUsers(users),
		}, nil
	}

	// もしメールが指定されていたら
	if req.Email != "" {
		// メールを検索する
		users, err := models.SearchUserByEmail(req.Email)
		if err != nil {
			return nil, err
		}

		return &SearchResult{
			Users: ModelUsersToUsers(users),
		}, nil
	}

	return nil, errors.New("not found")
}

func ModelUsersToUsers(modelUsers []models.User) []*User {
	var users []*User

	// モデルからユーザーを取得
	for _, user := range modelUsers {
		users = append(users, &User{
			UserID:        user.UserID,
			Name:          user.Name,
			Email:         user.Email,
			Labels:        ModelLabelsToLabels(user.Labels),
		})
	}

	return users
}

func ModelLabelsToLabels(modelLabels []models.Label) []*Label {
	var labels []*Label
	for _, label := range modelLabels {
		labels = append(labels, &Label{
			Name:  label.Name,
			Color: label.Color,
		})
	}
	return labels
}