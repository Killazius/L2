package main

import (
	"flag"
	"fmt"
	"l2/8/myntp"
	"os"
	"time"
)

var address = flag.String("address", "", "NTP server address")

func main() {
	flag.Parse()
	if *address == "" {
		fmt.Printf("no address provided, using default: %s\n", myntp.DefaultNTPServer)
		*address = myntp.DefaultNTPServer
	}
	t, err := myntp.GetTime(*address)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to get time from NTP:", err)
		os.Exit(1)

	}
	fmt.Println(t.Format(time.RFC1123))
}
