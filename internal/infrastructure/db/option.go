package db

type FindOption interface {
	apply(*option)
}

type option struct {
	query    []Query
	order    string
	offset   int
	limit    int
	preloads []string
	noOrder  bool
}

type optionFn func(*option)

func (f optionFn) apply(opt *option) {
	f(opt)
}

func WithQuery(query ...Query) FindOption {
	return optionFn(func(opt *option) {
		opt.query = query
	})
}

func WithOffset(offset int) FindOption {
	return optionFn(func(opt *option) {
		opt.offset = offset
	})
}

func WithLimit(limit int) FindOption {
	return optionFn(func(opt *option) {
		opt.limit = limit
	})
}

func WithOrder(order string) FindOption {
	return optionFn(func(opt *option) {
		opt.order = order
	})
}

func WithPreload(preloads []string) FindOption {
	return optionFn(func(opt *option) {
		opt.preloads = preloads
	})
}

// MATIKAN ORDER BY
func WithoutOrder() FindOption {
	return optionFn(func(opt *option) {
		opt.noOrder = true
	})
}

func getOption(opts ...FindOption) option {
	opt := option{
		query:    []Query{},
		offset:   0,
		limit:    1000,
		order:    "",       // default kosong, nanti di applyOptions pilih otomatis
		noOrder:  false,
	}

	for _, o := range opts {
		o.apply(&opt)
	}

	return opt
}
