package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

func init() {
	// Should be enabled manually. See:
	// https://docs.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#output-sequences
	err := enableTerminalSequences()

	if err != nil {
		fmt.Fprintln(
			os.Stderr,
			"Unable to enable terminal sequences.\n"+
				"Program output may have invalid formatting.",
		)
		fmt.Fprintln(os.Stderr, err)
	}
}

func enableTerminalSequences() error {
	fd := os.Stdout.Fd()
	handle := windows.Handle(fd)

	var mode uint32

	if err := windows.GetConsoleMode(handle, &mode); err != nil {
		return err
	}

	mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING

	if err := windows.SetConsoleMode(handle, mode); err != nil {
		return err
	}

	return nil
}
