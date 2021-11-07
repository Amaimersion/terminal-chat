package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Amaimersion/terminal-chat/chat"
)

func main() {
	flags, err := getChatFlags()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = chat.Run(flags)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getChatFlags() (chat.Flags, error) {
	var in, out, address, port string

	flag.StringVar(
		&in,
		"in",
		"/dev/stdin",
		"Where to read input from.",
	)
	flag.StringVar(
		&out,
		"out",
		"/dev/stdout",
		"Where to write output.",
	)
	flag.StringVar(
		&address,
		"address",
		"0.0.0.0",
		"What IP address to use for local server.",
	)
	flag.StringVar(
		&port,
		"port",
		"4444",
		"What TCP port to use for local server.",
	)

	flag.Parse()

	flags := chat.Flags{
		In:      nil,
		Out:     nil,
		Address: address,
		Port:    port,
	}

	if in == "/dev/stdin" {
		flags.In = os.Stdin
	} else {
		f, err := os.Open(in)

		if err != nil {
			return flags, err
		}

		flags.In = f
	}

	if out == "/dev/stdout" {
		flags.Out = os.Stdout
	} else {
		f, err := os.OpenFile(
			out,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND,
			os.ModePerm,
		)

		if err != nil {
			return flags, err
		}

		flags.Out = f
	}

	return flags, nil
}
