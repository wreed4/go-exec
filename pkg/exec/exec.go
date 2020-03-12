package exec

import (
	"bytes"
	"io/ioutil"
	"os"
	osexec "os/exec"
)

type Cmd struct {
	*osexec.Cmd
	_stdout string
	_stderr string
	_err    error
}

func Command(name string, args ...string) *Cmd {
	c := &Cmd{osexec.Command(name, args...), "", "", nil}
	return c
}

func (c Cmd) IsError() bool {
	return c._err != nil || c.ProcessState.ExitCode() != 0
}

func (c Cmd) Error() error {
	return c._err
}

func (c *Cmd) Dir(dir string) *Cmd {
	c.Cmd.Dir = dir
	return c
}

func (c *Cmd) Start() *Cmd {
	var out, err bytes.Buffer
	// do not capture if running in foreground
	if c.Stdout != os.Stdout {
		c.Stdout = &out
	}
	if c.Stderr != os.Stderr {
		c.Stderr = &err
	}
	c._err = c.Cmd.Start()
	return c
}

func (c *Cmd) Wait() *Cmd {
	c._err = c.Cmd.Wait()

	// do not capture if running in foreground
	if c.Stdout != os.Stdout {
		out := c.Stdout.(*bytes.Buffer)
		c._stdout = out.String()
	}
	if c.Stderr != os.Stderr {
		err := c.Stderr.(*bytes.Buffer)
		c._stderr = err.String()
	}
	return c
}

func (c *Cmd) Run() *Cmd {
	c = c.Start()
	if c.Error() != nil {
		return c
	}
	return c.Wait()
}

func (c *Cmd) RunFg() *Cmd {
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func (c *Cmd) Out() string { return c._stdout }
func (c *Cmd) Err() string { return c._stderr }

// func (c *Cmd) outputBuffer() (bytes.Buffer, bytes.Buffer) {
// 	var stderr bytes.Buffer
// 	var stdout bytes.Buffer
// 	c.Cmd.Stdout = &stdout
// 	c.Cmd.Stderr = &stderr
// 	c._err = c.Cmd.Run()

// 	return stdout, stderr
// }

// func (c *Cmd) OutputBytes() ([]byte, []byte, *Cmd) {
// 	stdout, stderr := c.outputBuffer()
// 	if c.IsError() {
// 		return nil, nil, c
// 	}

// 	return stdout.Bytes(), stderr.Bytes(), c
// }

// func (c *Cmd) Output() (string, string, *Cmd) {
// 	stdout, stderr := c.outputBuffer()
// 	if c.IsError() {
// 		return "", "", c
// 	}

// 	return stdout.String(), stderr.String(), c
// }

type BakedCmd func(...string) *Cmd

func (c *Cmd) Bake() BakedCmd {
	return func(args ...string) *Cmd {
		newc := &Cmd{&osexec.Cmd{
			Path:       c.Path,
			Args:       c.Args,
			Env:        c.Env,
			Dir:        c.Cmd.Dir,
			Stdin:      c.Stdin,
			Stdout:     c.Stdout,
			Stderr:     c.Stderr,
			ExtraFiles: c.ExtraFiles,
		}, "", "", nil}

		newc.Args = append(newc.Args, args...)
		return newc
	}
}

func WithTempdir(dir, pattern string, f func() error) error {
	name, err := ioutil.TempDir(dir, pattern)
	if err != nil {
		return err
	}
	defer os.Remove(name)
	dir, err = os.Getwd()
	if err != nil {
		return err
	}

	defer os.Chdir(dir)
	err = os.Chdir(name)
	if err != nil {
		return err
	}

	return f()
}
