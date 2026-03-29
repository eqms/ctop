package manager

type Runc struct{}

func NewRunc() *Runc {
	return &Runc{}
}

func (rc *Runc) Start() error {
	return ErrActionNotImpl
}

func (rc *Runc) Stop() error {
	return ErrActionNotImpl
}

func (rc *Runc) Remove() error {
	return ErrActionNotImpl
}

func (rc *Runc) Pause() error {
	return ErrActionNotImpl
}

func (rc *Runc) Unpause() error {
	return ErrActionNotImpl
}

func (rc *Runc) Restart() error {
	return ErrActionNotImpl
}

func (rc *Runc) Exec(cmd []string) error {
	return ErrActionNotImpl
}

func (rc *Runc) Commit(repo, tag string) error {
	return ErrActionNotImpl
}
