package wordmaker

import (
	"bytes"
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
			return nil, fmt.Errorf("Invalid config: %v", header)
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
	steps []chooser
}

func (p *Pattern) Run() chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, step := range p.steps {
			Debugf("step %q", step)
			choice, err := step.Choose()
			if err != nil {
				panic(err)
			}
			for _, c := range choice {
				ch <- string(c)
			}
		}
	}()

	return ch
}

func (p *Pattern) Choose() (string, error) {
	buf := bytes.Buffer{}
	for _, choice := range p.Items() {
		buf.WriteString(choice)
	}
	return buf.String(), nil
}

func (p *Pattern) Items() []string {
	items := []string{}
	for choice := range p.Run() {
		items = append(items, choice)
	}
	return items
}

func (p *Pattern) Append(ch chooser) error {
	p.steps = append(p.steps, ch)
	return nil
}

func (c *ChoiceList) Choose() (string, error) {
	var ch interface{}
	Debugf("Choose among %v", c.choices)
	ch, err := R.WeightedChoice(c.choices)
	if err != nil {
		return "", err
	}
	return ch.(R.Choice).Item.(chooser).Choose()
}

func (c *ChoiceList) Items() []string {
	out := []string{}
	for _, item := range c.choices {
		val, _ := item.Item.(chooser).Choose()
		out = append(out, val)
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
	Debugf("CHOICES %q", values)
	for _, v := range values {
		switch v.(type) {
		case string:
			nc := R.Choice{Weight: weight, Item: Choice{value: v.(string)}}
			cl.choices = append(cl.choices, nc)
		case []interface{}:
			nc := R.Choice{Weight: weight,
				Item: Choices("", v.([]interface{}), dropoff)}
			cl.choices = append(cl.choices, nc)
		case chooser:
			nc := R.Choice{Weight: weight, Item: v}
			cl.choices = append(cl.choices, nc)
		default:
			panic(fmt.Sprintf("Unknown choice type %T (%q)", v, v))
			// Debugf("WTF IS a %q (%T)", v, v)
		}
		weight = int(float64(weight) * dropoff)
	}
	return cl
}

func MakeChoices(name string, items chan item, dropoff float64) *ChoiceList {
	var pat *Pattern
	choices := []interface{}{}
	Debug("MakeChoices")

	// FIXME
	// within this loop, build up a sequence, append it to choices only when
	// reaching rightparen or end or slash
Loop:
	for i := range items {
		Debugf(" mc %v", i)
		switch i.typ {
		case itemChoice:
			Debugf("   mc a choice")
			if pat == nil {
				pat = NewPattern()
			}
			pat.Append(&Choice{value: i.val})
		case itemSlash:
			// append choice
			Debugf("  mc slash /")
			if pat != nil {
				choices = append(choices, pat)
				pat = nil
			}
		case itemLeftParen:
			// append choice
			Debugf("    mc ->")
			if pat != nil {
				choices = append(choices, pat)
				pat = nil
			}
			if pat == nil {
				pat = NewPattern()
			}
			pat.Append(MakeChoices("", items, dropoff))
		case itemRightParen:
			Debugf("    mc <-")
			break Loop
		}
	}

	// append choice if not nil
	if pat != nil {
		choices = append(choices, pat)
		pat = nil
	}
	return Choices(name, choices, dropoff)
}

func NewPattern() *Pattern {
	return &Pattern{steps: []chooser{}}
}

func MakePattern(items chan item, dropoff float64) *Pattern {
	Debug("MakePattern")
	pat := NewPattern()
Loop:
	for i := range items {
		Debugf("pat item %q", i)
		switch i.typ {
		case itemChoice:
			pat.steps = append(pat.steps, Choice{value: i.val})
		case itemLeftParen:
			pat.steps = append(pat.steps, MakeChoices("", items, dropoff))
		case itemRightParen:
			break Loop
		}
	}
	return pat
}
