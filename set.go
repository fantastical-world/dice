package dice

import (
	"fmt"
	"sort"
	"strings"
)

//Set holds custom dice that are backed by an expression.
//You can add dice to your set and roll them as often as needed.
type Set struct {
	Name string            `json:"name"`
	Dice map[string]string `json:"dice"`
}

//AddDice will store a roll expression as a custom dice in your set. The name provided can be passed to the RollDice function to roll the expression.
func (s *Set) AddDice(name string, expression string) error {
	if !ValidRollExpression(expression) {
		return fmt.Errorf("error %s is not a valid roll expression", expression)
	}

	if s.Dice == nil {
		s.Dice = make(map[string]string)
	}

	s.Dice[name] = expression

	return nil
}

//RollDice rolls the named custom dice's expression and returns its results.
func (s *Set) RollDice(name string) (rolls []int, sum int, err error) {
	if s.Dice == nil || len(s.Dice) == 0 {
		return rolls, sum, fmt.Errorf("error you do not have any dice in your set")
	}

	expression := s.Dice[name]

	if expression == "" {
		return rolls, sum, fmt.Errorf("error you do not have any dice named [%s] in your set", name)
	}

	rolls, sum, err = RollExpression(expression)

	return
}

//ListDice returns a string listing of all dice names and expressions in the set.
func (s *Set) ListDice() string {
	if s.Dice == nil || len(s.Dice) == 0 {
		return "no dice"
	}

	keys := make([]string, 0, len(s.Dice))
	for key := range s.Dice {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	sb := strings.Builder{}
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s %s\n", k, s.Dice[k]))
	}
	return sb.String()
}
