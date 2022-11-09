package workerpool

type Option func(*pool)

func WithBlock(block bool) Option {
	return func(p *pool) {
		p.block = block
	}
}

func WithPreAllocWorkers(preAlloc bool) Option {
	return func(p *pool) {
		p.preAlloc = preAlloc
	}
}
