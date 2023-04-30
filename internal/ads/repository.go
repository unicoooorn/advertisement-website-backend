package ads

import (
	"context"
)

type AdRepository interface {
	GetAdById(ctx context.Context, id int64) (*Ad, error)
	ListPublishedAds(ctx context.Context) ([]Ad, error)
	GetAllAds(ctx context.Context) ([]Ad, error)
	AddAd(ctx context.Context, ad Ad) (int64, error)
	UpdateById(ctx context.Context, id int64, ad Ad) error
	DeleteAdById(ctx context.Context, id int64) error
	DeleteAd(ctx context.Context, id int64)
}
