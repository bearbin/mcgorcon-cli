package main

import (
	"bufio"
	"code.google.com/p/gopass"
	"flag"
	"fmt"
	"github.com/mattn/go-isatty"
	"github.com/bearbin/mcgorcon"
	"io"
	"math/rand"
	"os"
	"time"
)

type configuration struct {
	Host     string
	Port     int
	Password string
}

func (c *configuration) Populate() {
	flag.StringVar(&c.Host, "host", "127.0.0.1", "the hostname of the server to connect to")
	flag.IntVar(&c.Port, "port", 25575, "the port the server is running on")
	flag.StringVar(&c.Password, "pass", "", "the password for the RCON service.")
	flag.Parse()
}

var config = configuration{}
var term bool

func init() {
	// Seed the RNG. Only needs doing once at startup.
	rand.Seed(time.Now().UTC().UnixNano())
	// Get the configuration from the available configuration methods.
	config.Populate()
	// Test if we have a terminal or a script as input.
	term = isatty.IsTerminal(os.Stdin.Fd())
}

func main() {
	if config.Password == "" {
		config.Password, _ = gopass.GetPass("Please enter the RCON server password: ")
	}
	client := mcgorcon.Dial(config.Host, config.Port, config.Password)
	stdin := bufio.NewReader(os.Stdin)
	for {
		if term {
			fmt.Print(">>> ")
		}
		input, err := stdin.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		} else if err != nil {
			panic(err)
		}
		output := client.SendCommand(input)
		fmt.Println(output)
	}

}
