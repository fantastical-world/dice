//Package dice offers functions and types for rolling any number of custom dice.
//You can use the various Roll functions directly by passing the number of dice to roll, the number of sides, modifiers, etc. The results returned will typically include the value of each dice rolled, their sum, and a modified sum if a modifer was provided.
//
//This package also provides a convenence function that excepts "Roll Expressions". A roll expression follows the typical 1d4+1 format used in most RPGs. With this single function you can satisfy most of your roll needs.
//
//Commonly used rolls can be saved for later use by creating a Set and adding them to it. A set acts like a dice bag except you can save expressions, not just dice.
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

//Roll rolls the specified number of n-sided dice and returns the rolled results and their sum.
func Roll(number int, sides int) (rolls []int, sum int) {
	rand.Seed(time.Now().UnixNano())
	rolls = make([]int, number)
	for i := 0; i < number; i++ {
		rolls[i] = rand.Intn(sides) + 1
		sum += rolls[i]
	}

	return
}

//RollAndModify rolls the specified number of n-sided dice then applies the provided modifier. The rolled results, their sum, and the modified sum will be returned.
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

//Modify applies the provided modifier to a roll and returns the value.
func Modify(rolledValue int, operator string, rollModifier int) (modifiedValue int) {
	modifiedValue = rolledValue
	switch operator {
	case "-":
		modifiedValue -= rollModifier
		return
	case "+":
		modifiedValue += rollModifier
		return
	}

	return
}

//ValidRollExpression validates that the provided expression is formatted correctly returning true if it is valid.
func ValidRollExpression(expression string) bool {
	return re.MatchString(expression)
}

//RollExpression will parse the provided roll expression and return its results.
func RollExpression(expression string) (rolls []int, sum int, err error) {
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
		return nil, 0, fmt.Errorf("not a valid roll expression, must be d# or #d# or #d#+# or #d#-# (e.g. d100, 1d4, 2d4+1, 2d6-2)")
	}

	match := re.FindStringSubmatch(expression)
	number, _ := strconv.Atoi(match[1])
	if number == 0 {
		number = 1
	}
	sides, _ := strconv.Atoi(match[2])

	//Handle min and max here
	if wantsMax {
		rolls, sum = RollMax(number, sides)
		if match[3] != "" {
			modifier, _ := strconv.Atoi(match[4])
			sum = Modify(sum, match[3], modifier)
		}
		return
	}

	if wantsMin {
		rolls, sum = RollMin(number, sides)
		if match[3] != "" {
			modifier, _ := strconv.Atoi(match[4])
			sum = Modify(sum, match[3], modifier)
		}
		return
	}

	//Anything after this comment uses the standard roll
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

	if dropLowest {
		lowest := rolls[0]
		for _, roll := range rolls {
			lowest = min(lowest, roll)
		}
		sum = sum - lowest
		return
	}

	if dropHighest {
		highest := rolls[0]
		for _, roll := range rolls {
			highest = max(highest, roll)
		}
		sum = sum - highest
		return
	}

	return
}

//RollMax rolls the specified number of n-sided dice then returns the rolled results and max value to use.
func RollMax(number int, sides int) (rolls []int, maxRoll int) {
	rolls, _ = Roll(number, sides)
	maxRoll = rolls[0]
	for _, roll := range rolls {
		maxRoll = max(maxRoll, roll)
	}

	return
}

//RollMin rolls the specified number of n-sided dice then returns the rolled results and min value to use.
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

//Version returns the current version.
func Version() string {
	return fmt.Sprintf("dice version: %s, build %s %s, %s\n", semVer, buildCommit, buildBranch, buildDate)
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
