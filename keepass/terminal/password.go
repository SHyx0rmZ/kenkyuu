package terminal

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SHyx0rmZ/kenkyuu/keepass"
	"golang.org/x/crypto/ssh/terminal"
)

func ReadPasswordFromTerminal() [32]byte {
	fd := int(os.Stdin.Fd())

	if !terminal.IsTerminal(fd) {
		bs, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalln(err)
		}
		return sha256.Sum256(bs)
	} else {
		fmt.Print("password: ")

		state, err := terminal.GetState(fd)
		if err != nil {
			log.Fatalln(err)
		}

		sigCh := make(chan os.Signal, 2)
		go func() {
			sig, ok := <-sigCh
			if !ok {
				return
			}

			err = terminal.Restore(fd, state)
			if err != nil {
				log.Println(err)
			}

			signal.Stop(sigCh)
			close(sigCh)
			for range sigCh {
			}

			proc, err := os.FindProcess(os.Getpid())
			if err != nil {
				log.Fatalln(err)
			}

			err = proc.Signal(sig)
			if err != nil {
				log.Fatalln(err)
			}
		}()
		signal.Notify(sigCh, os.Interrupt, os.Kill, syscall.SIGIO)

		bs, err := terminal.ReadPassword(fd)
		close(sigCh)
		fmt.Println()
		if err != nil {
			log.Fatalln(err)
		}

		return keepass.PasswordKey(bs)
	}
}
