package dicebag

import "fmt"

//Set hold a set of dice.
type Set struct {
	Name string
	Dice map[string]string
}

//AddDice adds a dice to your custom dice set.
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

//RollDice rolls the named dice.
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
