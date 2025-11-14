package pagination

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

type SortKey struct {
	Key   string
	Order string
}

type Pagination struct {
	enabled  bool
	limit    int64
	nextPage string
	page     int64
	sortKeys []SortKey

	Size int64
}

func (p Pagination) Enabled() bool {
	return p.enabled
}

func (p Pagination) Limit() int64 {
	return p.limit
}

func (p Pagination) NextPage() string {
	if p.Limit() < 1 {
		return ""
	}

	if p.Size < p.Limit() {
		return ""
	}

	return p.nextPage
}

func (p Pagination) Offset() int64 {
	if p.Limit() < 1 {
		return 0
	}

	if p.Page() <= 1 {
		return 0
	}

	return p.Limit() * (p.Page() - 1)
}

func (p Pagination) Page() int64 {
	if p.page < 1 {
		return 1
	}

	return p.page
}

func (p Pagination) SortKeys() []SortKey {
	// If no sort keys are provided, default to sorting by ID to ensure consistent results
	if p.sortKeys == nil {
		return []SortKey{}
	}

	return p.sortKeys
}

func (p Pagination) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey{}, &p)
}

func FromRequest(cfg Configuration, r *http.Request) (p Pagination, err error) {
	query := r.URL.Query()

	limit, err := getLimit(query.Get(cfg.LimitParameter), cfg.DefaultLimit, cfg.MinLimit, cfg.MaxLimit)
	if err != nil {
		return p, err
	}

	if limit < 0 {
		return p, nil
	}

	page, err := getPage(query.Get(cfg.PageParameter))
	if err != nil {
		return p, err
	}

	nextPage := getNextPage(cfg.PageParameter, page, r)

	sortKeys, err := getSorts(query[cfg.SortParameter], cfg.DefaultSort)
	if err != nil {
		return p, err
	}

	return Pagination{
		enabled:  true,
		limit:    limit,
		nextPage: nextPage,
		page:     page,
		sortKeys: sortKeys,
	}, nil
}

func getLimit(param string, def int64, min int64, max int64) (int64, error) {
	if param == "" {
		return def, nil
	}

	limit, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, err
	}

	if limit < min || limit > max {
		return 0, ErrInvalidLimit
	}

	return limit, nil
}

func getNextPage(key string, page int64, r *http.Request) string {
	nextPageQuery := r.URL.Query()
	nextPageQuery.Set(key, strconv.FormatInt(page+1, 10))

	nextPageURL := *r.URL
	nextPageURL.RawQuery = nextPageQuery.Encode()

	return nextPageURL.String()
}

func getPage(param string) (int64, error) {
	if param == "" {
		return 1, nil
	}

	page, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, err
	}
	if page < 1 {
		return 0, ErrInvalidPage
	}

	return page, nil
}

func getSorts(params []string, def string) ([]SortKey, error) {
	// If no sort keys are provided, default to sorting by ID to ensure consistent results
	if len(params) == 0 {
		key, order, _ := strings.Cut(def, ",")

		return []SortKey{
			{
				Key:   key,
				Order: order,
			},
		}, nil
	}

	sortKeys := make([]SortKey, len(params))

	for i, val := range params {
		key, order, _ := strings.Cut(val, ",")

		if key == "" {
			return nil, ErrMissingSortKey
		}

		sortKeys[i] = SortKey{Key: key, Order: order}
	}

	return sortKeys, nil
}
