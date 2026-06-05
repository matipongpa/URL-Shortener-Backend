package shortener

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/matipongpa/url-shortener/storage"
)

type Storage interface {
	Save(ctx context.Context, code, url string) error
	Get(ctx context.Context, code string) (string, error)
}

type Shortener struct{ store Storage }

func New(store Storage) *Shortener {
	return &Shortener{store: store} // this is trivial — just write it
}

func (s *Shortener) Shorten(ctx context.Context, url string) (string, error) {
	var err error
	var code string
	for range 5 {
		code, err = generateShortCode(6)
		if err != nil {
			return "", err
		}
		err = s.store.Save(ctx, code, url)
		if err == nil {
			return code, nil
		}
		if !errors.Is(err, storage.ErrCodeExists) {
			return "", err
		}
	}
	return "", errors.New("shortener: failed to generate unique code after 5 attempts")
}

func (s *Shortener) Resolve(ctx context.Context, code string) (string, error) {
	return s.store.Get(ctx, code)
}

func generateShortCode(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes)[:length], nil
}
