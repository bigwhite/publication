package workerpool

type Option func(*Pool)

func WithBlock(block bool) Option {
	return func(p *Pool) {
		p.block = block
	}
}

func WithPreAllocWorkers(preAlloc bool) Option {
	return func(p *Pool) {
		p.preAlloc = preAlloc
	}
}
