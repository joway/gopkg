package zeromalloc

var _ Allocator = (*sallocator)(nil)

// safe allocator
type sallocator struct {
	Allocator
	mu uint32
}

func NewSafe(unit int, page int, limit int) (Allocator, error) {
	a, err := NewUnsafe(unit, page, limit)
	if err != nil {
		return nil, err
	}
	return &sallocator{
		Allocator: a,
	}, nil
}

func (a *sallocator) Alloc() (p uintptr, err error) {
	lock(&a.mu)
	p, err = a.Allocator.Alloc()
	unlock(&a.mu)
	return p, err
}

func (a *sallocator) Free(p uintptr) {
	lock(&a.mu)
	a.Allocator.Free(p)
	unlock(&a.mu)
}

func (a *sallocator) Close() (err error) {
	lock(&a.mu)
	err = a.Close()
	unlock(&a.mu)
	return err
}
