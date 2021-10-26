//Package dice offers functions and types for rolling any number of custom dice.
//You can use the various Roll functions directly by passing the number of dice
//to roll, the number of sides, modifiers, etc. The results returned will typically
//include the value of each dice rolled, their sum, and a modified sum if a modifer was provided.
//
//This package also provides a convenence function that excepts "Roll Expressions".
//A roll expression follows the typical 1d4+1 format used in most RPGs.
//With this single function you can satisfy most of your roll needs.
//
//Commonly used rolls can be saved for later use by creating a Set and adding them to it.
//A set acts like a dice bag except you can save expressions, not just dice.
package dice

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	RollExpressionRE               = regexp.MustCompile(`^([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?((\+|-)([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?)?$`)         //entire string is a roll expression (e.g. "2d6+3") pre-pair-regex (`^([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?$`)
	ContainsRollExpressionRE       = regexp.MustCompile(`\s*([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?((\+|-)([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?)?\s*`)     //any roll expression in a string (e.g. "Hi roll {{2d6+3}} to hit.") pre-pair-regex (`\s*([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?\s*`)
	ContainsRollExpressionBracedRE = regexp.MustCompile(`{{\s*([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?((\+|-)([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?)?\s*}}`) //same as above, but will include braces in matches, pre-pair-regex (`{{\s*([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?\s*}}`)
)

//Roll rolls the specified number of n-sided dice and returns the rolled results and their sum.
func Roll(number int, sides int) ([]int, int, error) {
	if number < 0 {
		return nil, 0, ErrInvalidNumberOfDice
	}
	if sides < 0 {
		return nil, 0, ErrInvalidNumberOfSides
	}
	rand.Seed(time.Now().UnixNano())
	rolls := make([]int, number)
	sum := 0
	for i := 0; i < number; i++ {
		rolls[i] = rand.Intn(sides) + 1
		sum += rolls[i]
	}

	return rolls, sum, nil
}

//RollAndModify rolls the specified number of n-sided dice then applies the provided modifier.
//The rolled results, their sum, and the modified sum will be returned.
//An error is returned if the operator is anything other than + or -.
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

//Modify applies the provided modifier to a roll and returns the value.
//An error is returned if the operator is anything other than + or -.
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

//ValidRollExpression validates that the provided expression is formatted correctly returning true if it is valid.
func ValidRollExpression(expression string) bool {
	return RollExpressionRE.MatchString(expression)
}

//ContainsValidRollExpression checks the provided string for valid roll expressions and returns count of valid found.
func ContainsValidRollExpression(data string) int {
	all := ContainsRollExpressionRE.FindAllString(data, -1)
	if all == nil {
		return 0
	}

	return len(all)
}

//RollExpression will parse the provided roll expression and return its results.
//An error is returned if the expression is invalid. The min: and max: can cause
//an error if they are used with an expression pair.
func RollExpression(expression string) ([]int, int, error) {
	//check for a special prefix
	var wantsMax, wantsMin, halfResult, doubleResult, dropLowest, dropHighest bool
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

	if strings.HasPrefix(expression, "dropL:") {
		dropLowest = true
		expression = strings.ReplaceAll(expression, "dropL:", "")
	}

	if strings.HasPrefix(expression, "dropH:") {
		dropHighest = true
		expression = strings.ReplaceAll(expression, "dropH:", "")
	}

	//simple 1d4+1 style (#d#+|-# or #d# or d#)
	if !ValidRollExpression(expression) {
		return nil, 0, ErrInvalidRollExpression
	}

	hasSecondExpression := false
	match := RollExpressionRE.FindStringSubmatch(expression)
	if match[5] != "" {
		hasSecondExpression = true
	}

	//min: and max: prefix is not valid if expression is a pair/double expression.
	if hasSecondExpression && (wantsMax || wantsMin) {
		return nil, 0, ErrInvalidRollExpression
	}

	number, _ := strconv.Atoi(match[1])
	//convert the absence of a number to mean 1 to satisfy d6 like shorthand, otherwise it was a 0
	if number == 0 && match[1] == "" {
		number = 1
	}
	sides, _ := strconv.Atoi(match[2])

	//handle min and max here
	if wantsMax {
		rolls, sum, _ := RollMax(number, sides)
		if match[3] != "" {
			modifier, _ := strconv.Atoi(match[4])
			sum, _ = Modify(sum, match[3], modifier)
		}
		return rolls, sum, nil
	}

	if wantsMin {
		rolls, sum, _ := RollMin(number, sides)
		if match[3] != "" {
			modifier, _ := strconv.Atoi(match[4])
			sum, _ = Modify(sum, match[3], modifier)
		}
		return rolls, sum, nil
	}

	//anything after this comment uses the standard roll or roll modify
	var rolls []int
	var sum int
	if match[3] == "" {
		rolls, sum, _ = Roll(number, sides)
	} else {
		modifier, _ := strconv.Atoi(match[4])
		rolls, _, sum, _ = RollAndModify(number, sides, match[3], modifier)
	}

	//the drop prefixes only apply to the first expression, the second expression is treated like a modifier
	if dropLowest {
		lowest := 0
		for r, roll := range rolls {
			if r == 0 {
				lowest = roll
				continue
			}
			lowest = min(lowest, roll)
		}
		sum = sum - lowest
	}

	if dropHighest {
		highest := 0
		for r, roll := range rolls {
			if r == 0 {
				highest = roll
				continue
			}
			highest = max(highest, roll)
		}
		sum = sum - highest
	}

	//let's handle second expression if provided
	if hasSecondExpression {
		var secondRolls []int
		var secondSum int

		secondNumber, _ := strconv.Atoi(match[7])
		//convert the absence of a number to mean 1 to satisfy d6 like shorthand, otherwise it was a 0
		if secondNumber == 0 && match[7] == "" {
			secondNumber = 1
		}
		secondSides, _ := strconv.Atoi(match[8])
		if match[9] == "" {
			secondRolls, secondSum, _ = Roll(secondNumber, secondSides)
		} else {
			secondModifier, _ := strconv.Atoi(match[10])
			secondRolls, _, secondSum, _ = RollAndModify(secondNumber, secondSides, match[9], secondModifier)
		}

		switch match[6] {
		case "-":
			sum -= secondSum
		case "+":
			sum += secondSum
		}

		rolls = append(rolls, secondRolls...)
	}

	if halfResult {
		sum = sum / 2
		return rolls, sum, nil
	}

	if doubleResult {
		sum = sum * 2
		return rolls, sum, nil
	}

	return rolls, sum, nil
}

//RollMax rolls the specified number of n-sided dice then returns the rolled results and max value to use.
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

//RollMin rolls the specified number of n-sided dice then returns the rolled results and min value to use.
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

//RollChallenge rolls an expression against a provided value. The rolled value must be greater
//than the challenge value to succeed. If desired the challenge can succeed on equal values
//by setting equalSucceeds to true. You can also be alerted when specific values are rolled
//by providing a slice of values, if any were rolled they will be returned.
//
//An error is returned if the expression is not a valid roll expression.
func RollChallenge(expression string, against int, equalSucceeds bool, alertOn []int) (bool, int, []int, error) {
	rolls, result, err := RollExpression(expression)
	if err != nil {
		return false, 0, nil, err
	}

	succeeded := (result > against)

	if !succeeded && equalSucceeds {
		succeeded = (result == against)
	}

	var found []int
	if len(alertOn) > 0 {
		for _, roll := range rolls {
			for _, check := range alertOn {
				if roll == check {
					found = append(found, check)
					break
				}
			}
		}
	}

	return succeeded, result, found, nil
}

func RollString(value string) string {
	rolledValue := value
	if !ContainsRollExpressionBracedRE.MatchString(value) {
		return value
	}

	match := ContainsRollExpressionBracedRE.FindAllStringSubmatch(value, 99) //limit to 99 rolls per value
	for _, m := range match {
		expression := strings.ReplaceAll(m[0], "{{", "")
		expression = strings.ReplaceAll(expression, "}}", "")
		_, sum, _ := RollExpression(strings.Trim(expression, " "))
		rolledValue = strings.Replace(rolledValue, m[0], strconv.Itoa(sum), 1)
	}

	return rolledValue
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
