//Package dice offers functions for rolling any number of custom dice. Roll functions typically return the result of each dice rolled, their sum, and a modified sum if a modifer was provided.
//This package also uses "Roll Expressions" that are typically used in RPGs. Roll expressions will be translated into the appropriate number of n-sided dice, and any modifiers.
package dice

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	buildBranch string
	buildCommit string
	buildDate   string
	semVer      string
	re          = regexp.MustCompile(`^([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?$`)
)

//Roll rolls the specified number of n-sided dice and returns the rolled results and the sum of those rolls.
func Roll(number int, sides int) (rolls []int, sum int) {
	rand.Seed(time.Now().UnixNano())
	rolls = make([]int, number)
	for i := 0; i < number; i++ {
		rolls[i] = rand.Intn(sides) + 1
		sum += rolls[i]
	}

	return
}

//RollAndModify rolls the specified number of n-sided dice then adjusts the sum with modifier.
func RollAndModify(number int, sides int, operator string, rollModifier int) (rolls []int, sum int, modifiedSum int) {
	rolls, sum = Roll(number, sides)
	modifiedSum = sum

	switch operator {
	case "-":
		modifiedSum -= rollModifier
		return
	case "+":
		modifiedSum += rollModifier
		return
	}

	return
}

//ValidRollExpression checks the provided expression and returns true if valid.
func ValidRollExpression(expression string) bool {
	return re.MatchString(expression)
}

//RollExpression roll the expression in simple 1d4+1 style (#d#+|-# or #d#).
func RollExpression(expression string) (rolls []int, sum int, err error) {
	//check for a special prefix
	var wantsMax, wantsMin, halfResult, doubleResult bool
	if strings.HasPrefix(expression, "max:") {
		wantsMax = true
		expression = strings.ReplaceAll(expression, "max:", "")
	}

	if strings.HasPrefix(expression, "min:") {
		wantsMin = true
		expression = strings.ReplaceAll(expression, "min:", "")
	}

	if strings.HasPrefix(expression, "half:") {
		halfResult = true
		expression = strings.ReplaceAll(expression, "half:", "")
	}

	if strings.HasPrefix(expression, "dub:") {
		doubleResult = true
		expression = strings.ReplaceAll(expression, "dub:", "")
	}

	//simple 1d4+1 style (#d#+|-# or #d# or d#)
	if !ValidRollExpression(expression) {
		return nil, 0, fmt.Errorf("not a valid roll expression, must be d# or #d# or #d#+# or #d#-# (e.g. d100, 1d4, 2d4+1, 2d6-2)")
	}

	match := re.FindStringSubmatch(expression)
	number, _ := strconv.Atoi(match[1])
	if number == 0 {
		number = 1
	}
	sides, _ := strconv.Atoi(match[2])

	if wantsMax {
		rolls, sum = RollMax(number, sides)
		return
	}

	if wantsMin {
		rolls, sum = RollMin(number, sides)
		return
	}

	if match[3] == "" {
		rolls, sum = Roll(number, sides)
	} else {
		modifier, _ := strconv.Atoi(match[4])
		rolls, _, sum = RollAndModify(number, sides, match[3], modifier)
	}

	if halfResult {
		sum = sum / 2
		return
	}

	if doubleResult {
		sum = sum * 2
		return
	}

	return
}

//RollMax returns rolls and max value to use.
func RollMax(number int, sides int) (rolls []int, maxRoll int) {
	rolls, _ = Roll(number, sides)
	maxRoll = rolls[0]
	for _, roll := range rolls {
		maxRoll = max(maxRoll, roll)
	}

	return
}

//RollMin returns rolls and min value to use.
func RollMin(number int, sides int) (rolls []int, minRoll int) {
	rolls, _ = Roll(number, sides)
	minRoll = rolls[0]
	for _, roll := range rolls {
		minRoll = min(minRoll, roll)
	}

	return
}

//RollChallenge rolls the expression against a provided value. The value must be greater than the challenge value to succeed.
//If desired the challenge can succeed on equal values by setting equalSucceeds to true.
//An error is returned if the expression is not a valid roll expression.
func RollChallenge(expression string, against int, equalSucceeds bool, alert []int) (succeeded bool, result int, found []int, err error) {
	rolls, result, err := RollExpression(expression)
	if err != nil {
		return false, 0, nil, err
	}

	succeeded = (result > against)

	if !succeeded && equalSucceeds {
		succeeded = (result == against)
	}

	if alert != nil && (len(alert) > 0) {
		for _, roll := range rolls {
			for _, check := range alert {
				if roll == check {
					found = append(found, check)
					break
				}
			}
		}
	}

	return
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
