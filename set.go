package dice

import (
	"fmt"
	"sort"
	"sync"
)

//Set holds custom dice that are backed by an expression.
//You can add dice to your set and roll them as often as needed.
type Set struct {
	m    sync.RWMutex
	Name string            `json:"name"`
	Dice map[string]string `json:"dice"`
}

//AddDice will store a roll expression as a custom dice in your set. The name provided can be passed to the RollDice function to roll the expression.
func (s *Set) AddDice(name string, expression string) error {
	s.m.Lock()
	defer s.m.Unlock()
	if !ValidRollExpression(expression) {
		return ErrInvalidRollExpression
	}

	if s.Dice == nil {
		s.Dice = make(map[string]string)
	}

	s.Dice[name] = expression

	return nil
}

//RemoveDice will remove the roll expression saved under the name provided.
func (s *Set) RemoveDice(name string) {
	s.m.Lock()
	defer s.m.Unlock()

	delete(s.Dice, name)
}

//RollDice rolls the named custom dice's expression and returns its results.
func (s *Set) RollDice(name string) (rolls []int, sum int, err error) {
	s.m.RLock()
	defer s.m.RUnlock()
	if s.Dice == nil || len(s.Dice) == 0 {
		return rolls, sum, ErrEmptyDiceSet
	}

	expression := s.Dice[name]

	if expression == "" {
		return rolls, sum, ErrDiceNotFound
	}

	rolls, sum, err = RollExpression(expression)

	return
}

//ListDice returns a listing of all dice names and expressions in the set.
func (s *Set) ListDice() []string {
	s.m.RLock()
	defer s.m.RUnlock()
	if s.Dice == nil || len(s.Dice) == 0 {
		return nil
	}

	keys := make([]string, 0, len(s.Dice))
	for key := range s.Dice {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var list []string
	for _, k := range keys {
		list = append(list, fmt.Sprintf("%s,%s", k, s.Dice[k]))
	}

	return list
}
