package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	W "github.com/jpellerin/wordmaker/wordmaker"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "mkwords"
	app.Usage = "Make some fake words using awkwords-like rules"

	app.Flags = []cli.Flag{
		cli.IntFlag{"words", 100, "number of words to generate"},
		cli.Float64Flag{"dropoff", 0.7, "rate of dropoff in choice lists"},
	}
	app.Action = func(c *cli.Context) {
		rulefile := "rules.aw"
		words := c.Int("words")
		dropoff := c.Float64("dropoff")
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
			fmt.Printf("%v ", w)
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
