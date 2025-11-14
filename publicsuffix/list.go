package publicsuffix

import (
	"context"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/weppos/publicsuffix-go/publicsuffix"
)

type List struct {
	cache    *lru.TwoQueueCache[string, DomainName]
	errCache *lru.TwoQueueCache[string, error]
	psl      *publicsuffix.List
}

func NewList(path string, cacheSize int, errCacheSize int) (List, error) {
	list := List{}

	if cacheSize > 0 {
		cache, err := lru.New2Q[string, DomainName](cacheSize)
		if err != nil {
			return list, err
		}

		list.cache = cache
	}

	if errCacheSize > 0 {
		errCache, err := lru.New2Q[string, error](errCacheSize)
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
	if l.errCache != nil {
		if err, cached := l.errCache.Get(name); cached {
			return DomainName{}, err
		}
	}

	if l.cache != nil {
		if res, cached := l.cache.Get(name); cached {
			return res, nil
		}
	}

	dn, err := publicsuffix.ParseFromListWithOptions(l.psl, name, DefaultFindOptions)
	if err != nil {
		if l.errCache != nil {
			l.errCache.Add(name, err)
		}

		return DomainName{}, err
	}

	fqdn, domain, tld, subdomain := Expand(dn)

	res := DomainName{
		FQDN:      fqdn,
		Domain:    domain,
		Subdomain: subdomain,
		TLD:       tld,
	}

	if l.cache != nil {
		l.cache.Add(name, res)
	}

	return res, nil
}

func (l List) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey{}, l)
}
