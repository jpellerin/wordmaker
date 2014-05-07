package wordmaker

import (
	R "github.com/jmcvetta/randutil"
)

type chooser interface {
	Choose() (string, error)
}

type Choice struct {
	value string
}

type ChoiceList struct {
	choices []R.Choice
}

func (c *ChoiceList) Choose() (string, error) {
	var ch interface{}
	ch, err := R.WeightedChoice(c.choices)
	if err != nil {
		return "", err
	}
	return ch.(R.Choice).Item.(chooser).Choose()
}

func (c Choice) Choose() (string, error) {
	return c.value, nil
}

func Choices(values []interface{}, dropoff float64) *ChoiceList {
	cl := &ChoiceList{}
	weight := 1000
	for _, v := range values {
		switch v.(type) {
		case string:
			nc := R.Choice{Weight: weight, Item: Choice{value: v.(string)}}
			cl.choices = append(cl.choices, nc)
		case []interface{}:
			nc := R.Choice{Weight: weight,
				Item: Choices(v.([]interface{}), dropoff)}
			cl.choices = append(cl.choices, nc)
		}
		weight = int(float64(weight) * dropoff)
	}
	return cl
}
