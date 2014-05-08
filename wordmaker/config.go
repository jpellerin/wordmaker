package wordmaker

import (
	// "fmt"
	R "github.com/jmcvetta/randutil"
	"strings"
)

type Config struct {
	Name     string
	classes  map[string]*ChoiceList
	patterns []*Pattern
}

func NewConfig(name string) *Config {
	return &Config{
		Name:     name,
		classes:  map[string]*ChoiceList{},
		patterns: []*Pattern{},
	}
}

func (c *Config) AddChoiceClass(cls *ChoiceList) error {
	c.classes[cls.Name] = cls
	return nil
}

func (c *Config) AddPattern(pat *Pattern) error {
	c.patterns = append(c.patterns, pat)
	return nil
}

func (c *Config) Class(name string) *ChoiceList {
	return c.classes[name]
}

func (c *Config) Word() (string, error) {
	word := []string{}
	pat, err := pick(c.patterns)
	if err != nil {
		return "", err
	}

	for val := range pat.Run() {
		cls := c.Class(val)
		if cls == nil {
			word = append(word, val)
		} else {
			val, err := cls.Choose()
			if err != nil {
				return "", err
			}
			word = append(word, val)
		}
	}
	return strings.Join(word, ""), nil
}

func pick(pats []*Pattern) (*Pattern, error) {
	var winner *Pattern
	length := len(pats)
	i, err := R.IntRange(0, length)
	winner = pats[i]
	return winner, err
}
