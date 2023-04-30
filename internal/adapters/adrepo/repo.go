package adrepo

import (
	"context"
	"errors"
	"homework9/internal/ads"
	"sync"
)

type RepositoryMap struct {
	repo   map[int64]ads.Ad
	lastId int64
	mx     *sync.RWMutex
}

func New() *RepositoryMap {
	return &RepositoryMap{repo: make(map[int64]ads.Ad), lastId: -1, mx: &sync.RWMutex{}}
}

// - создание нового объявления
//- публикация или снятие объявления с публикации
//- изменение текста объявления

func (r *RepositoryMap) GetAdById(ctx context.Context, id int64) (*ads.Ad, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	ad, ok := r.repo[id]
	if !ok {
		return nil, errors.New("ad not found")
	}
	return &ad, nil
}

func (r *RepositoryMap) AddAd(ctx context.Context, ad ads.Ad) (int64, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.lastId++
	id := r.lastId
	r.repo[id] = ad
	return id, nil
}

func (r *RepositoryMap) DeleteAdById(ctx context.Context, id int64) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	_, ok := r.repo[id]
	if !ok {
		return errors.New("ad not found")
	}
	delete(r.repo, id)
	return nil
}

func (r *RepositoryMap) UpdateById(ctx context.Context, id int64, ad ads.Ad) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	_, ok := r.repo[id]
	if !ok {
		return errors.New("ad not found")
	}
	r.repo[id] = ad
	return nil
}

func (r *RepositoryMap) ListPublishedAds(ctx context.Context) ([]ads.Ad, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	res := make([]ads.Ad, 0)
	for _, ad := range r.repo {
		if ad.Published {
			res = append(res, ad)
		}
	}
	return res, nil
}

func (r *RepositoryMap) GetAllAds(ctx context.Context) ([]ads.Ad, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	res := make([]ads.Ad, 0)
	for _, ad := range r.repo {
		res = append(res, ad)
	}
	return res, nil
}

func (r *RepositoryMap) DeleteAd(ctx context.Context, id int64) {
	r.mx.Lock()
	defer r.mx.Unlock()
	delete(r.repo, id)
}
