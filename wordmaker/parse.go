package wordmaker

import (
	"fmt"
	R "github.com/jmcvetta/randutil"
)

func Parse(name string, input []string, dropoff float64) (*Config, error) {
	cfg := NewConfig(name)
	for _, line := range input {
		_, items := Lex(name, line)
		header := <-items
		switch header.typ {
		case itemClass:
			if err := cfg.AddChoiceClass(MakeChoices(header.val, items, dropoff)); err != nil {
				return nil, err
			}
		case itemPattern:
			if err := cfg.AddPattern(MakePattern(items, dropoff)); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("Invalid config")
		}
	}
	return cfg, nil
}

type chooser interface {
	Choose() (string, error)
	Items() []string
}

type Choice struct {
	value string
}

type ChoiceList struct {
	Name    string
	choices []R.Choice
}

type Pattern struct {
	steps []R.Choice
}

func (p *Pattern) Run() chan string {
	ch := make(chan string)
	go func() {
		for _, step := range p.steps {
			choice, err := step.Item.(chooser).Choose()
			if err != nil {
				panic(err)
			}
			for _, c := range choice {
				ch <- string(c)
			}
		}
		defer close(ch)
	}()

	return ch
}

func (c *ChoiceList) Choose() (string, error) {
	var ch interface{}
	ch, err := R.WeightedChoice(c.choices)
	if err != nil {
		return "", err
	}
	return ch.(R.Choice).Item.(chooser).Choose()
}

func (c *ChoiceList) Items() []string {
	out := []string{}
	for _, item := range c.choices {
		out = append(out, item.Item.(Choice).value)
	}
	return out
}

func (c Choice) Choose() (string, error) {
	return c.value, nil
}

func (c Choice) Items() []string {
	return []string{c.value}
}

func Choices(name string, values []interface{}, dropoff float64) *ChoiceList {
	cl := &ChoiceList{Name: name}
	weight := 1000
	for _, v := range values {
		switch v.(type) {
		case string:
			nc := R.Choice{Weight: weight, Item: Choice{value: v.(string)}}
			cl.choices = append(cl.choices, nc)
		case []interface{}:
			nc := R.Choice{Weight: weight,
				Item: Choices("", v.([]interface{}), dropoff)}
			cl.choices = append(cl.choices, nc)
		}
		weight = int(float64(weight) * dropoff)
	}
	return cl
}

func MakeChoices(name string, items chan item, dropoff float64) *ChoiceList {
	choices := []interface{}{}
	// fmt.Print("\nMakeChoices\n")
Loop:
	for i := range items {
		// fmt.Printf(" mc %v\n", i)
		switch i.typ {
		case itemChoice:
			// fmt.Printf("   mc a choice\n")
			choices = append(choices, i.val)
		case itemLeftParen:
			// fmt.Printf("    mc ->\n")
			choices = append(choices, MakeChoices("", items, dropoff))
		case itemRightParen:
			// fmt.Printf("    mc <-\n")
			break Loop
		}
	}

	return Choices(name, choices, dropoff)
}

func MakePattern(items chan item, dropoff float64) *Pattern {
	pat := &Pattern{steps: []R.Choice{}}
	// fmt.Print("\nMakePattern\n")
	for i := range items {
		// fmt.Printf("pat item %q\n", i)
		switch i.typ {
		case itemChoice:
			pat.steps = append(pat.steps, R.Choice{Item: Choice{value: i.val}})
		case itemLeftParen:
			pat.steps = append(pat.steps, R.Choice{Item: MakeChoices("", items, dropoff)})
		}
	}
	return pat
}
