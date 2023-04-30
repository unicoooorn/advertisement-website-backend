package tests

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	grpcPort "homework9/internal/ports/grpc"
)

func TestGRRPCCreateUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService()
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Name)
}

func TestGRRPCCreateAd(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService()
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Vladimir"})
	assert.NoError(t, err, "client.CreateUser")
	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Amerike konec", Text: "Ya otvechayu", UserId: user.Id})
	assert.NoError(t, err, "client.CreateAd")

	assert.Equal(t, "Amerike konec", res.Title)
	assert.Equal(t, "Ya otvechayu", res.Text)
	assert.Equal(t, int64(0), res.Id)
}

func TestGRRPCChangeAdStatus(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService()
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Vladimir"})
	assert.NoError(t, err, "client.CreateUser")
	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Amerike konec", Text: "Ya otvechayu", UserId: user.Id})
	assert.NoError(t, err, "client.CreateAd")
	adUpd, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: ad.Id, UserId: user.Id, Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")

	assert.True(t, adUpd.Published)
}

func TestGRRPCUpdateAd(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService()
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Vladimir"})
	assert.NoError(t, err, "client.CreateUser")
	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Amerike konec", Text: "Ya otvechayu", UserId: user.Id})
	assert.NoError(t, err, "client.CreateAd")
	adUpd, err := client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{AdId: ad.Id, UserId: user.Id, Title: "Ya poshutil", Text: "A vi poverili?"})
	assert.NoError(t, err, "client.UpdateAd")

	assert.Equal(t, adUpd.Text, "A vi poverili?")
	assert.Equal(t, adUpd.Title, "Ya poshutil")
	assert.Equal(t, ad.Id, adUpd.Id)
}

func TestGRRPCListAds(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService()
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user1, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Vladimir"})
	assert.NoError(t, err, "client.CreateUser")
	user2, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Dmitry"})
	assert.NoError(t, err, "client.CreateUser")
	ad1, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Amerike konec", Text: "Ya otvechayu", UserId: user1.Id})
	assert.NoError(t, err, "client.CreateAd")
	ad2, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Europe tozhe", Text: "Ya guarantee", UserId: user1.Id})
	assert.NoError(t, err, "client.CreateAd")
	ad3, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Nam horosho", Text: "Da", UserId: user2.Id})
	assert.NoError(t, err, "client.CreateAd")
	ad1, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: ad1.Id, UserId: user1.Id, Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")
	ad2, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: ad2.Id, UserId: user1.Id, Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")
	ad3, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: ad3.Id, UserId: user2.Id, Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")
	ads, err := client.ListAds(ctx, &emptypb.Empty{})
	assert.NoError(t, err, "client.ListAds")

	assert.Contains(t, []string{ad1.Title, ad2.Title, ad3.Title}, ads.List[0].Title)
	assert.Contains(t, []string{ad1.Title, ad2.Title, ad3.Title}, ads.List[1].Title)
	assert.Contains(t, []string{ad1.Title, ad2.Title, ad3.Title}, ads.List[2].Title)
}
func TestGRRPCDeleteUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService()
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	_, err = client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Vladimir"})
	assert.NoError(t, err, "client.CreateUser")
	user1, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Dmitry"})
	assert.NoError(t, err, "client.CreateUser")
	_, err = client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{Id: user1.Id})
	assert.NoError(t, err, "client.DeleteUser")
	user, err := client.GetUser(ctx, &grpcPort.GetUserRequest{Id: user1.Id})
	assert.Error(t, err, "not found")
	assert.Nil(t, user)
}

func TestGRRPCDeleteAd(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService()
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Vladimir"})
	assert.NoError(t, err, "client.CreateUser")
	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "What", Text: "The", UserId: u.Id})
	assert.NoError(t, err, "client.CreateAd")
	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: u.Id, Published: true, AdId: ad.Id})
	assert.NoError(t, err, "client.ChangeAdStatus")
	_, err = client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AdId: ad.Id})
	assert.NoError(t, err, "client.DeleteAd")
	ads, err := client.ListAds(ctx, &emptypb.Empty{})
	assert.NoError(t, err, "client.ListAds")
	assert.NotContains(t, ads.List, ad)
}

func TestGRRPCGetUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService()
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	_, err = client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Dmitry"})
	assert.NoError(t, err, "client.GetUser")
	userCreated, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Pavel"})
	assert.NoError(t, err, "client.CreateUser")
	userGot, err := client.GetUser(ctx, &grpcPort.GetUserRequest{Id: userCreated.Id})
	assert.NoError(t, err, "client.GetUser")
	assert.Equal(t, userCreated.Name, userGot.Name)
}
