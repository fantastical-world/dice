package dicebag

import (
	"fmt"
	"sort"
	"strings"
)

//Set hold a set of dice.
type Set struct {
	Name string
	dice map[string]string
}

//AddDice adds a dice to your custom dice set.
func (s *Set) AddDice(name string, expression string) error {
	if !ValidRollExpression(expression) {
		return fmt.Errorf("error %s is not a valid roll expression", expression)
	}

	if s.dice == nil {
		s.dice = make(map[string]string)
	}

	s.dice[name] = expression

	return nil
}

//RollDice rolls the named custom dice.
func (s *Set) RollDice(name string) (rolls []int, sum int, err error) {
	if s.dice == nil || len(s.dice) == 0 {
		return rolls, sum, fmt.Errorf("error you do not have any dice in your set")
	}

	expression := s.dice[name]

	if expression == "" {
		return rolls, sum, fmt.Errorf("error you do not have any dice named [%s] in your set", name)
	}

	rolls, sum, err = RollExpression(expression)

	return
}

//ListDice returns a string listing of all dice in set.
func (s *Set) ListDice() string {
	if s.dice == nil || len(s.dice) == 0 {
		return "no dice"
	}

	keys := make([]string, 0, len(s.dice))
	for key := range s.dice {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	sb := strings.Builder{}
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s %s\n", k, s.dice[k]))
	}
	return sb.String()
}
