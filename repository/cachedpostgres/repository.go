package cachedpostgres

import "github.com/go-url-shortener/shortener"

type cachedRepository struct {
	repository      shortener.RedirectRepository
	redisRepository shortener.RedirectRepository
}

func NewCachedRepository(postgresRepository shortener.RedirectRepository, redisRepository shortener.RedirectRepository) shortener.RedirectRepository {
	return &cachedRepository{
		postgresRepository,
		redisRepository,
	}
}

func (r *cachedRepository) Find(code string) (*shortener.Redirect, error) {
	redirect, err := r.redisRepository.Find(code)

	if err != nil {
		redirect, err := r.repository.Find(code)
		if err == nil {
			r.redisRepository.Store(redirect)
		}
		return redirect, err
	}

	return redirect, nil
}

func (r *cachedRepository) Store(redirect *shortener.Redirect) error {
	return r.repository.Store(redirect)
}
