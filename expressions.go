package dice

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	RollExpressionRE               = regexp.MustCompile(`^([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?((\+|-)([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?)?$`)         //entire string is a roll expression (e.g. "2d6+3") pre-pair-regex (`^([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?$`)
	ContainsRollExpressionRE       = regexp.MustCompile(`\s*([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?((\+|-)([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?)?\s*`)     //any roll expression in a string (e.g. "Hi roll {{2d6+3}} to hit.") pre-pair-regex (`\s*([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?\s*`)
	ContainsRollExpressionBracedRE = regexp.MustCompile(`{{\s*([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?((\+|-)([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?)?\s*}}`) //same as above, but will include braces in matches, pre-pair-regex (`{{\s*([0-9]*)[d]([0-9]+)(\+|-)?([0-9]+)?\s*}}`)
)

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
