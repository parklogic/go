package publicsuffix

import (
	"context"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/weppos/publicsuffix-go/publicsuffix"
)

type lookupKey struct {
	name   string
	noStar bool
}

type List struct {
	cache    *lru.TwoQueueCache[lookupKey, DomainName]
	errCache *lru.TwoQueueCache[lookupKey, error]
	psl      *publicsuffix.List
}

func NewList(path string, cacheSize int, errCacheSize int) (List, error) {
	list := List{}

	if cacheSize > 0 {
		cache, err := lru.New2Q[lookupKey, DomainName](cacheSize)
		if err != nil {
			return list, err
		}

		list.cache = cache
	}

	if errCacheSize > 0 {
		errCache, err := lru.New2Q[lookupKey, error](errCacheSize)
		if err != nil {
			return list, err
		}

		list.errCache = errCache
	}

	if path == "" {
		path = DefaultListPath
	}

	psl, err := publicsuffix.NewListFromFile(path, DefaultParserOptions)
	if err != nil {
		return list, err
	}

	list.psl = psl

	return list, nil
}

func (l List) Parse(name string) (DomainName, error) {
	return l.parse(name, false)
}

func (l List) ParseNoStar(name string) (DomainName, error) {
	return l.parse(name, true)
}

func (l List) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey{}, l)
}

func (l List) parse(name string, noStar bool) (DomainName, error) {
	key := lookupKey{name: name, noStar: noStar}

	if res, ok, err := l.cached(key); ok {
		return res, err
	}

	parsed, err := publicsuffix.ParseFromListWithOptions(l.psl, name, findOptions(noStar))
	if err != nil {
		l.rememberErr(key, err)

		return DomainName{}, err
	}

	err = normalise(parsed)
	if err != nil {
		l.rememberErr(key, err)

		return DomainName{}, err
	}

	res := expand(parsed)
	if l.cache != nil {
		l.cache.Add(key, res)
	}

	return res, nil
}

func (l List) cached(key lookupKey) (DomainName, bool, error) {
	var zero DomainName

	if l.errCache != nil {
		err, hit := l.errCache.Get(key)
		if hit {
			return zero, true, err //nolint:wrapcheck // the cache stores the error parse already returned
		}
	}

	if l.cache != nil {
		res, hit := l.cache.Get(key)
		if hit {
			return res, true, nil
		}
	}

	return zero, false, nil
}

func (l List) rememberErr(key lookupKey, err error) {
	if l.errCache != nil {
		l.errCache.Add(key, err)
	}
}
