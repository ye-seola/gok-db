package sendmsg

import (
	"bufio"
	"fmt"
	"gokdb/internal/utils"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

var (
	cmdMutex sync.Mutex
	started                 = false
	stdin    io.WriteCloser = nil
	stdout   *bufio.Scanner = nil
)

func mustStart() {
	var err error

	cmd := exec.Command("/system/bin/app_process", "/", "SendMsg")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("CLASSPATH=%s", filepath.Join(
		utils.GetExecDir(), "sendmsg.dex",
	)))

	stdin, err = cmd.StdinPipe()
	if err != nil {
		panic(fmt.Errorf("failed to create stdin pipe: %w", err))
	}

	_stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(fmt.Errorf("failed to create stdout pipe: %w", err))
	}
	stdout = bufio.NewScanner(_stdout)

	if err := cmd.Start(); err != nil {
		panic(fmt.Errorf("failed to start command: %w", err))
	}

	started = true
}
