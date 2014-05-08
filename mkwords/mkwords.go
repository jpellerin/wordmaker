package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	W "github.com/jpellerin/wordmaker/wordmaker"
	"log"
	"os"
)

func init() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} [rulefile] [options] command

   The default rulefile is "rules.aw"

VERSION:
   {{.Version}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}

RULE FILE FORMAT:
   * Lines starting with an uppercase character define word parts
   * Lines starting with r define word patterns
   * Slashes (/) separate options, which are in order of descending weight
   * Parentheses may be used to group a set of options

   Example:

   V:a/i/u/ei/ao/ia/ai
   C:p/t/k/s/m/n/b/w/x/y/ts/l/sh/ch
   T:p/t/k
   F:s
   N:m/n/rn/rl/nd/ng
   r:CV(N/-CV/-CVN/)
`
}

func main() {
	app := cli.NewApp()
	app.Name = "mkwords"
	app.Usage = "Make some fake words using awkwords-like rules"
	app.Version = "0.1.1"

	app.Flags = []cli.Flag{
		cli.IntFlag{"words", 100, "number of words to generate"},
		cli.Float64Flag{"dropoff", 0.7, "rate of dropoff in choice lists"},
		cli.BoolFlag{"one-per-line", "output one word per line"},
	}
	app.Action = func(c *cli.Context) {
		rulefile := "rules.aw"
		words := c.Int("words")
		dropoff := c.Float64("dropoff")
		sep := " "
		if c.Bool("one-per-line") {
			sep = "\n"
		}
		if len(c.Args()) > 0 {
			rulefile = c.Args()[0]
		}
		input, err := readlines(rulefile)
		if err != nil {
			log.Fatalf("Error reading rules file %v: %v", rulefile, err)
		}
		cfg, err := W.Parse("mkwords", input, dropoff)
		if err != nil {
			log.Fatalf("Failed to parse rules file %v: %v", rulefile, err)
		}
		for i := 0; i < words; i++ {
			w, err := cfg.Word()
			if err != nil {
				log.Fatalf("Word generation failure: %v", err)
			}
			fmt.Printf("%v", w)
			fmt.Print(sep)
		}
		fmt.Print("\n")
	}
	app.Run(os.Args)
}

func readlines(filename string) ([]string, error) {
	lines := []string{}
	ff, err := os.Open(filename)
	if err != nil {
		return lines, err
	}
	defer ff.Close()
	scanner := bufio.NewScanner(ff)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
