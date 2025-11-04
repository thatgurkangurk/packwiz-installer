package core

import (
	"context"
	"net/http"
	"net/url"
	"sync"

	"golang.org/x/sync/errgroup"
)

type RepoOptFn func(r *Repository)

type Repository struct {
	Url            *url.URL
	Pack           *PackToml
	Index          *IndexToml
	Metafiles      []*MetafileToml
	PackHashFormat string
	PackHash       string
	httpClient     *http.Client
}

func NewRepository(url *url.URL, hashFormat, hash string) *Repository {
	return &Repository{
		Url:            url,
		PackHashFormat: hashFormat,
		PackHash:       hash,
		httpClient:     http.DefaultClient,
	}
}

func (r *Repository) loadPack(ctx context.Context) (*PackToml, error) {
	var (
		data []byte
		err  error
	)

	if r.PackHash == "" {
		data, err = httpGetBytes(ctx, r.httpClient, r.Url.String())
		if err != nil {
			return nil, err
		}
	} else {
		data, err = httpGetValidBytes(ctx, r.httpClient, r.Url.String(), r.PackHashFormat, r.PackHash)
		if err != nil {
			return nil, err
		}
	}

	pack, err := parsePackToml(data)
	if err != nil {
		return nil, err
	}
	r.Pack = pack
	return pack, nil
}

func (r *Repository) loadIndex(ctx context.Context) (*IndexToml, error) {
	if r.Pack == nil {
		_, err := r.loadPack(ctx)
		if err != nil {
			return nil, err
		}
	}

	data, err := httpGetValidBytes(
		ctx,
		r.httpClient,
		r.IndexUrl().String(),
		r.Pack.Index.HashFormat,
		r.Pack.Index.Hash,
	)
	if err != nil {
		return nil, err
	}

	index, err := parseIndexToml(data)
	if err != nil {
		return nil, err
	}
	r.Index = index
	return index, nil
}

func (r *Repository) loadMetafiles(ctx context.Context) ([]*MetafileToml, error) {
	if r.Index == nil {
		_, err := r.loadIndex(ctx)
		if err != nil {
			return nil, err
		}
	}

	var mods = make([]*MetafileToml, 0, len(r.Index.Files))
	eg := errgroup.Group{}
	mutex := sync.Mutex{}
	for _, file := range r.Index.Files {
		indexedFile := file
		eg.Go(func() error {
			if !indexedFile.Metafile {
				return nil
			}

			metafileUrl := r.IndexUrl().JoinPath("..", indexedFile.File).String()
			hashFmt := indexedFile.HashFormat
			if hashFmt == "" {
				hashFmt = r.Index.HashFormat
			}
			data, err := httpGetValidBytes(ctx, r.httpClient, metafileUrl, hashFmt, indexedFile.Hash)
			if err != nil {
				return err
			}

			mod, err := parseMetafileToml(data)
			if err != nil {
				return err
			}
			mod.IndexName = indexedFile.File

			mutex.Lock()
			mods = append(mods, mod)
			mutex.Unlock()
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	r.Metafiles = mods
	return mods, nil
}

func (r *Repository) Load(ctx context.Context) error {
	_, err := r.loadMetafiles(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) BaseUrl() *url.URL {
	return r.Url.JoinPath("..")
}

func (r *Repository) IndexUrl() *url.URL {
	return r.BaseUrl().JoinPath(r.Pack.Index.File)
}
