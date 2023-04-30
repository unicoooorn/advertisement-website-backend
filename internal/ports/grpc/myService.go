package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/adapters/userrepo"
	"homework9/internal/ads"
	"homework9/internal/app"
	"homework9/internal/users"
)

type AdService struct {
	app app.App
}

func NewService() AdServiceServer {
	return AdService{
		app: app.NewApp(adrepo.New(), userrepo.New()).(app.MyApp),
	}
}

func (as AdService) CreateAd(ctx context.Context, reqBody *CreateAdRequest) (*AdResponse, error) {
	_, err := as.app.GetUserByID(ctx, reqBody.UserId)
	if err == app.ErrNotFound {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	a, err := as.app.CreateAd(ctx, reqBody.Title, reqBody.Text, reqBody.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return newAdResponse(a), nil
}
func (as AdService) ChangeAdStatus(ctx context.Context, reqBody *ChangeAdStatusRequest) (*AdResponse, error) {
	_, err := as.app.GetUserByID(ctx, reqBody.UserId)
	if err == app.ErrNotFound {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	a, err := as.app.UpdateStatusById(ctx, reqBody.AdId, reqBody.Published, reqBody.UserId)
	if err == app.ErrAccessDenied {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return newAdResponse(a), nil
}
func (as AdService) UpdateAd(ctx context.Context, in *UpdateAdRequest) (*AdResponse, error) {
	_, err := as.app.GetUserByID(ctx, in.UserId)
	if err == app.ErrNotFound {
		return nil, status.Error(codes.NotFound, "User not found")
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	a, err := as.app.UpdateAdById(ctx, in.AdId, in.Title, in.Text, in.UserId)
	if err == app.ErrAccessDenied {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return newAdResponse(a), nil
}
func (as AdService) ListAds(ctx context.Context, in *emptypb.Empty) (*ListAdResponse, error) {
	l, err := as.app.ListPublishedAds(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return newMultipleAdsResponse(l), nil
}
func (as AdService) CreateUser(ctx context.Context, in *CreateUserRequest) (*UserResponse, error) {
	u, err := as.app.CreateUser(ctx, in.Name, "")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return newUserResponse(u), nil
}
func (as AdService) GetUser(ctx context.Context, in *GetUserRequest) (*UserResponse, error) {
	u, err := as.app.GetUserByID(ctx, in.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return newUserResponse(u), nil
}
func (as AdService) DeleteUser(ctx context.Context, in *DeleteUserRequest) (*emptypb.Empty, error) {
	err := as.app.DeleteUser(ctx, in.Id)
	if err == app.ErrNotFound {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
func (as AdService) DeleteAd(ctx context.Context, in *DeleteAdRequest) (*emptypb.Empty, error) {
	err := as.app.DeleteAd(ctx, in.AdId, in.AuthorId)
	if err == app.ErrNotFound {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err == app.ErrAccessDenied {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func newAdResponse(ad *ads.Ad) *AdResponse {
	return &AdResponse{Title: ad.Title,
		Text:      ad.Text,
		Id:        ad.ID,
		AuthorId:  ad.AuthorID,
		Published: ad.Published}
}

func newUserResponse(u *users.User) *UserResponse {
	return &UserResponse{
		Id:   u.ID,
		Name: u.Nickname,
	}
}

func newMultipleAdsResponse(ads []ads.Ad) *ListAdResponse {
	res := make([]*AdResponse, 0)
	for _, ad := range ads {
		res = append(res, newAdResponse(&ad))
	}
	return &ListAdResponse{
		List: res,
	}
}
