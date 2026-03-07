package manager

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	api "github.com/fsouza/go-dockerclient"
	"golang.org/x/term"
)

type Docker struct {
	id     string
	client *api.Client
}

func NewDocker(client *api.Client, id string) *Docker {
	return &Docker{
		id:     id,
		client: client,
	}
}

func (dc *Docker) Exec(cmd []string) error {
	execCmd, err := dc.client.CreateExec(api.CreateExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
		Container:    dc.id,
		Tty:          true,
	})
	if err != nil {
		return err
	}

	// Set host terminal to raw mode for interactive shell
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("cannot set terminal raw mode: %v", err)
	}
	defer term.Restore(fd, oldState)

	// Set initial TTY size
	if w, h, err := term.GetSize(fd); err == nil {
		dc.client.ResizeExecTTY(execCmd.ID, h, w)
	}

	// Handle terminal resize (SIGWINCH)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)
	defer signal.Stop(sigCh)
	go func() {
		for range sigCh {
			if w, h, err := term.GetSize(fd); err == nil {
				dc.client.ResizeExecTTY(execCmd.ID, h, w)
			}
		}
	}()

	// With Tty: true, Docker sends raw terminal data (no frame headers),
	// so we write directly to os.Stdout instead of using a frameWriter.
	cw, err := dc.client.StartExecNonBlocking(execCmd.ID, api.StartExecOptions{
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Tty:          true,
		RawTerminal:  true,
	})
	if err != nil {
		return err
	}
	if cw != nil {
		return cw.Wait()
	}
	return nil
}

func (dc *Docker) Start() error {
	c, err := dc.client.InspectContainer(dc.id)
	if err != nil {
		return fmt.Errorf("cannot inspect container: %v", err)
	}

	if err := dc.client.StartContainer(c.ID, c.HostConfig); err != nil {
		return fmt.Errorf("cannot start container: %v", err)
	}
	return nil
}

func (dc *Docker) Stop() error {
	if err := dc.client.StopContainer(dc.id, 3); err != nil {
		return fmt.Errorf("cannot stop container: %v", err)
	}
	return nil
}

func (dc *Docker) Remove() error {
	if err := dc.client.RemoveContainer(api.RemoveContainerOptions{ID: dc.id}); err != nil {
		return fmt.Errorf("cannot remove container: %v", err)
	}
	return nil
}

func (dc *Docker) Pause() error {
	if err := dc.client.PauseContainer(dc.id); err != nil {
		return fmt.Errorf("cannot pause container: %v", err)
	}
	return nil
}

func (dc *Docker) Unpause() error {
	if err := dc.client.UnpauseContainer(dc.id); err != nil {
		return fmt.Errorf("cannot unpause container: %v", err)
	}
	return nil
}

func (dc *Docker) Restart() error {
	if err := dc.client.RestartContainer(dc.id, 3); err != nil {
		return fmt.Errorf("cannot restart container: %v", err)
	}
	return nil
}

func (dc *Docker) Commit(repo, tag string) error {
	_, err := dc.client.CommitContainer(api.CommitContainerOptions{
		Container:  dc.id,
		Repository: repo,
		Tag:        tag,
	})
	if err != nil {
		return fmt.Errorf("cannot commit container: %v", err)
	}
	return nil
}
