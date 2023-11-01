// Package dice offers functions and types for rolling any number of custom dice.
// You can use the various Roll functions directly by passing the number of dice
// to roll, the number of sides, modifiers, etc. The results returned will typically
// include the value of each dice rolled, their sum, and a modified sum if a modifier was provided.
//
// This package also provides a convenience function that excepts "Roll Expressions".
// A roll expression follows the typical 1d4+1 format used in most RPGs.
// With this single function you can satisfy most of your roll needs.
//
// Commonly used rolls can be saved for later use by creating a Set and adding them to it.
// A set acts like a die bag except you can save expressions, not just dice.
package dice

import (
	"time"
)

// Roll rolls the specified number of n-sided dice and returns the rolled results and their sum.
func Roll(number int, sides int) ([]int, int, error) {
	if number < 0 {
		return nil, 0, ErrInvalidNumberOfDice
	}
	if sides < 0 {
		return nil, 0, ErrInvalidNumberOfSides
	}
	s := New(time.Now().UnixNano())
	rolls := s.RandomNRange(number, 1, sides, false)
	sum := 0
	for _, roll := range rolls {
		sum += roll
	}

	return rolls, sum, nil
}

// RollAndModify rolls the specified number of n-sided dice then applies the provided modifier.
// The rolled results, their sum, and the modified sum will be returned.
// An error is returned if the operator is anything other than + or -.
func RollAndModify(number int, sides int, operator string, rollModifier int) ([]int, int, int, error) {
	if number < 0 {
		return nil, 0, 0, ErrInvalidNumberOfDice
	}
	if sides < 0 {
		return nil, 0, 0, ErrInvalidNumberOfSides
	}
	rolls, sum, _ := Roll(number, sides)
	modifiedSum := sum

	switch operator {
	case "-":
		modifiedSum -= rollModifier
	case "+":
		modifiedSum += rollModifier
	default:
		return nil, 0, 0, ErrInvalidOperator
	}

	return rolls, sum, modifiedSum, nil
}

// Modify applies the provided modifier to a roll and returns the value.
// An error is returned if the operator is anything other than + or -.
func Modify(rolledValue int, operator string, rollModifier int) (int, error) {
	modifiedValue := rolledValue
	switch operator {
	case "-":
		modifiedValue -= rollModifier
	case "+":
		modifiedValue += rollModifier
	default:
		return 0, ErrInvalidOperator
	}

	return modifiedValue, nil
}

// RollMax rolls the specified number of n-sided dice then returns the rolled results and max value to use.
func RollMax(number int, sides int) ([]int, int, error) {
	if number < 0 {
		return nil, 0, ErrInvalidNumberOfDice
	}
	if sides < 0 {
		return nil, 0, ErrInvalidNumberOfSides
	}
	rolls, _, _ := Roll(number, sides)
	maxRoll := 0
	for r, roll := range rolls {
		if r == 0 {
			maxRoll = roll
			continue
		}
		maxRoll = max(maxRoll, roll)
	}

	return rolls, maxRoll, nil
}

// RollMin rolls the specified number of n-sided dice then returns the rolled results and min value to use.
func RollMin(number int, sides int) ([]int, int, error) {
	if number < 0 {
		return nil, 0, ErrInvalidNumberOfDice
	}
	if sides < 0 {
		return nil, 0, ErrInvalidNumberOfSides
	}
	rolls, _, _ := Roll(number, sides)
	minRoll := 0
	for r, roll := range rolls {
		if r == 0 {
			minRoll = roll
			continue
		}
		minRoll = min(minRoll, roll)
	}

	return rolls, minRoll, nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
