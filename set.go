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
	dice map[string]string
}

//AddDice will store a roll expression as a custom dice in your set. The name provided can be passed to the RollDice function to roll the expression.
func (s *Set) AddDice(name string, expression string) error {
	s.m.Lock()
	defer s.m.Unlock()
	if !ValidRollExpression(expression) {
		return ErrInvalidRollExpression
	}

	if s.dice == nil {
		s.dice = make(map[string]string)
	}

	s.dice[name] = expression

	return nil
}

//RemoveDice will remove the roll expression saved under the name provided.
func (s *Set) RemoveDice(name string) {
	s.m.Lock()
	defer s.m.Unlock()

	delete(s.dice, name)
}

//RollDice rolls the named custom dice's expression and returns its results.
func (s *Set) RollDice(name string) (rolls []int, sum int, err error) {
	s.m.RLock()
	defer s.m.RUnlock()
	if s.dice == nil || len(s.dice) == 0 {
		return rolls, sum, ErrEmptyDiceSet
	}

	expression := s.dice[name]

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
	if s.dice == nil || len(s.dice) == 0 {
		return nil
	}

	keys := make([]string, 0, len(s.dice))
	for key := range s.dice {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var list []string
	for _, k := range keys {
		list = append(list, fmt.Sprintf("%s,%s", k, s.dice[k]))
	}

	return list
}
