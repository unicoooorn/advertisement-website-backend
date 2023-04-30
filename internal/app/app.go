package app

import (
	"context"
	"errors"
	"homework9/internal/adapters/userrepo"
	"homework9/internal/ads"
	"homework9/internal/users"
	"time"
)

type FilterOpts struct {
	ID           int64
	Title        string
	AuthorID     int64
	ModifiedTime time.Time
	CreatedTime  time.Time
	Hidden       bool
}

type App interface {
	CreateAd(ctx context.Context, title string, text string, authorId int64) (*ads.Ad, error)
	UpdateStatusById(ctx context.Context, id int64, status bool, authorId int64) (*ads.Ad, error)
	UpdateAdById(ctx context.Context, id int64, title string, text string, authorId int64) (*ads.Ad, error)
	ListPublishedAds(ctx context.Context) ([]ads.Ad, error)
	GetAdsByFilter(ctx context.Context, opts FilterOpts) ([]ads.Ad, error)
	DeleteAd(ctx context.Context, id int64, userId int64) error

	GetUserByID(ctx context.Context, userID int64) (*users.User, error)
	CreateUser(ctx context.Context, nickname string, email string) (*users.User, error)
	UpdateUserByID(ctx context.Context, updatedID int64, nickname string, email string, updaterID int64) (*users.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

var ErrAccessDenied = errors.New("forbidden")
var ErrNotFound = errors.New("not found")

type MyApp struct {
	adRepository   ads.AdRepository
	userRepository users.UserRepository
}

func NewApp(adRepo ads.AdRepository, userRepo users.UserRepository) App {
	return MyApp{adRepository: adRepo, userRepository: userRepo}
}

func (m MyApp) CreateAd(ctx context.Context, title string, text string, authorId int64) (*ads.Ad, error) {
	a := ads.Ad{Title: title, Text: text, AuthorID: authorId, Published: false, Created: time.Now(), Modified: time.Now()}
	id, err := m.adRepository.AddAd(ctx, a)
	if err != nil {
		return nil, err
	}

	a.ID = id

	return &a, nil
}

func (m MyApp) UpdateStatusById(ctx context.Context, id int64, status bool, authorId int64) (*ads.Ad, error) {
	a, err := m.adRepository.GetAdById(ctx, id)
	if err != nil {
		return nil, err
	}
	if a.AuthorID != authorId {
		return nil, ErrAccessDenied
	}
	changed := ads.Ad{
		Title:     a.Title,
		Text:      a.Text,
		AuthorID:  a.AuthorID,
		Published: status,
		Created:   a.Created,
	}
	err = m.adRepository.UpdateById(ctx, id, changed)
	if err != nil {
		return nil, err
	}
	return &changed, nil
}

func (m MyApp) UpdateAdById(ctx context.Context, id int64, title string, text string, authorId int64) (*ads.Ad, error) {
	a, err := m.adRepository.GetAdById(ctx, id)
	if err != nil {
		return nil, err
	}
	if a.AuthorID != authorId {
		return nil, ErrAccessDenied
	}
	changed := ads.Ad{
		Title:     title,
		Text:      text,
		AuthorID:  a.AuthorID,
		Published: a.Published,
		Created:   a.Created,
		Modified:  time.Now(),
	}
	err = m.adRepository.UpdateById(ctx, id, changed)
	if err != nil {
		return nil, err
	}
	return &changed, nil
}

func (m MyApp) ListPublishedAds(ctx context.Context) ([]ads.Ad, error) {
	res, err := m.adRepository.ListPublishedAds(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m MyApp) CreateUser(ctx context.Context, nickname string, email string) (*users.User, error) {
	u := users.User{Nickname: nickname, Email: email}
	id, err := m.userRepository.AddUser(ctx, u)
	if err != nil {
		return nil, err
	}
	u.ID = id
	return &u, nil
}

func (m MyApp) UpdateUserByID(ctx context.Context, updatedID int64, nickname string, email string, updaterID int64) (*users.User, error) {
	u, err := m.userRepository.GetUserByID(ctx, updatedID)
	if err != nil {
		return nil, err
	}
	if u.ID != updaterID {
		return nil, ErrAccessDenied
	}
	changed := users.User{Nickname: nickname, Email: email}
	err = m.userRepository.UpdateByID(ctx, updatedID, changed)
	if err != nil {
		return nil, err
	}
	return &changed, nil
}

func (m MyApp) GetAdsByFilter(ctx context.Context, opts FilterOpts) ([]ads.Ad, error) {
	adsAll, err := m.adRepository.GetAllAds(ctx)
	if err != nil {
		return nil, err
	}
	adsFiltered := make([]ads.Ad, 0)
	for _, ad := range adsAll {
		if opts.ID != 0 && ad.ID != opts.ID {
			continue
		}
		if opts.Title != "" && ad.Title != opts.Title {
			continue
		}
		if opts.AuthorID != 0 && ad.AuthorID != opts.AuthorID {
			continue
		}
		if !opts.CreatedTime.IsZero() && !ad.Created.Equal(opts.CreatedTime) {
			continue
		}
		if !opts.ModifiedTime.IsZero() && !ad.Modified.Equal(opts.ModifiedTime) {
			continue
		}
		adsFiltered = append(adsFiltered, ad)
	}
	return adsFiltered, nil
}

func (m MyApp) GetUserByID(ctx context.Context, userID int64) (*users.User, error) {
	u, err := m.userRepository.GetUserByID(ctx, userID)
	if err == userrepo.ErrNotFound {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return u, nil
}

func (m MyApp) DeleteAd(ctx context.Context, id int64, userId int64) error {
	ad, err := m.adRepository.GetAdById(ctx, id)
	if err != nil {
		return err
	}
	if ad.AuthorID != userId {
		return ErrAccessDenied
	}
	if err == ErrNotFound {
		return err
	}
	m.adRepository.DeleteAd(ctx, id)
	return nil
}

func (m MyApp) DeleteUser(ctx context.Context, id int64) error {
	_, err := m.userRepository.GetUserByID(ctx, id)
	if err == ErrNotFound {
		return err
	}
	m.userRepository.DeleteUser(ctx, id)
	return nil
}
