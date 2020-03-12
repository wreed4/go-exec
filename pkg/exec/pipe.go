package exec

import (
	"fmt"
	"os"
)

type PipedCmd struct {
	In *PipedCmd
	*Cmd
	final bool
}

func Pipe(cmds ...*Cmd) *PipedCmd {
	var p = &PipedCmd{}
	if len(cmds) < 2 {
		p._err = fmt.Errorf("pipe requires at least two commands")
		return p
	}

	p.Cmd = cmds[0]
	return p.Pipe(cmds[1:]...)
}

func (p *PipedCmd) Pipe(cmds ...*Cmd) *PipedCmd {
	for _, c := range cmds {
		outpipe, err := p.StdoutPipe()
		if err != nil {
			p._err = err
			return p
		}

		c.Stdin = outpipe

		p.final = false
		p = &PipedCmd{
			In:    p,
			Cmd:   c,
			final: true,
		}
	}

	return p
}

func (p *PipedCmd) String() string {
	var s string
	if p.In != nil {
		s = p.In.String() + " | "
	}
	return s + p.Cmd.String()
}

func (p *PipedCmd) Start() *Cmd {
	if p.In != nil {
		in := p.In.Start()
		if in.Error() != nil {
			return in
		}
	}
	if p.final {
		return p.Cmd.Start()
	}
	p.Cmd._err = p.Cmd.Cmd.Start()
	return p.Cmd
}

func (p *PipedCmd) Wait() *Cmd {
	if p.In != nil {
		in := p.In.Wait()
		if in.Error() != nil {
			return in
		}
	}
	if p.final {
		return p.Cmd.Wait()
	}
	p.Cmd._err = p.Cmd.Cmd.Wait()
	return p.Cmd
}

func (p *PipedCmd) Run() *Cmd {
	c := p.Start()
	if c.Error() != nil {
		return c
	}
	return p.Wait()
}

func (p *PipedCmd) RunFg() *Cmd {
	p.Cmd.Stdout = os.Stdout
	p.Cmd.Stderr = os.Stderr
	return p.Run()
}
