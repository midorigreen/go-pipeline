package pipeline

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

func OutputStderr(stdErr io.Writer, commands ...[]string) ([]byte, error) {
	cmds := make([]*exec.Cmd, len(commands))
	var err error

	for i, c := range commands {
		cmds[i] = exec.Command(c[0], c[1:]...)
		if i > 0 {
			if cmds[i].Stdin, err = cmds[i-1].StdoutPipe(); err != nil {
				return nil, err
			}
		}
		cmds[i].Stderr = stdErr
	}
	var out bytes.Buffer
	cmds[len(cmds)-1].Stdout = &out
	for _, c := range cmds {
		if err = c.Start(); err != nil {
			return nil, err
		}
	}
	for _, c := range cmds {
		if err = c.Wait(); err != nil {
			return nil, err
		}
	}
	return out.Bytes(), nil
}

func Output(commands ...[]string) ([]byte, error) {
	return OutputStderr(os.Stderr, commands...)
}
