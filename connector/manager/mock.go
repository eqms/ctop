package manager

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) Start() error {
	return ErrActionNotImpl
}

func (m *Mock) Stop() error {
	return ErrActionNotImpl
}

func (m *Mock) Remove() error {
	return ErrActionNotImpl
}

func (m *Mock) Pause() error {
	return ErrActionNotImpl
}

func (m *Mock) Unpause() error {
	return ErrActionNotImpl
}

func (m *Mock) Restart() error {
	return ErrActionNotImpl
}

func (m *Mock) Exec(cmd []string) error {
	return ErrActionNotImpl
}

func (m *Mock) Commit(repo, tag string) error {
	return ErrActionNotImpl
}
