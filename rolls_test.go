package dice

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRollingASingleDice(t *testing.T) {
	//roll a single six sided dice
	rolls, sum := Roll(1, 6)
	//only one roll should be present
	if len(rolls) != 1 {
		t.Errorf("expected 1 roll result, but found %d results\n", len(rolls))
	}
	//single roll means sum and roll should be equal
	if rolls[0] != sum {
		t.Errorf("expected roll result and sum to be equal, but roll was %d and sum was %d\n", rolls[0], sum)
	}
	//roll should be 1, 2, 3, 4, 5, or 6
	if rolls[0] < 1 || rolls[0] > 6 {
		t.Errorf("expected roll result to be 1, 2, 3, 4, 5, or 6, but roll was %d\n", rolls[0])
	}
}

func TestRollingMultipleDice(t *testing.T) {
	//roll four six sided dice
	rolls, sum := Roll(4, 6)
	//four rolls should be present
	if len(rolls) != 4 {
		t.Errorf("expected 4 roll results, but found %d results\n", len(rolls))
	}
	//check sum
	expectedSum := 0
	for _, roll := range rolls {
		expectedSum += roll
	}
	if expectedSum != sum {
		t.Errorf("expected roll results to equal sum when added, but roll results equaled %d and sum was %d\n", expectedSum, sum)
	}
	//rolls should be 1, 2, 3, 4, 5, or 6
	for _, roll := range rolls {
		if roll < 1 || roll > 6 {
			t.Errorf("expected roll result to be 1, 2, 3, 4, 5, or 6, but roll was %d\n", roll)
		}
	}
}

func TestRollWithAddModifier(t *testing.T) {
	//roll two six sided dice and add three to result sum
	rolls, sum, modSum := RollAndModify(2, 6, "+", 3)
	if len(rolls) != 2 {
		t.Errorf("expected 2 roll results, but found %d results\n", len(rolls))
	}
	//check sum
	expectedSum := 0
	for _, roll := range rolls {
		expectedSum += roll
	}
	if expectedSum != sum {
		t.Errorf("expected roll results to equal sum when added, but roll results equaled %d and sum was %d\n", expectedSum, sum)
	}
	//check modified sum, it should be sum plus three
	if modSum != (sum + 3) {
		t.Errorf("expected modified sum to be sum + 3, but modified sum equaled %d and sum was %d\n", modSum, sum)
	}
}

func TestRollWithSubtractModifier(t *testing.T) {
	//roll two ten sided dice and subtract four to result sum
	rolls, sum, modSum := RollAndModify(2, 10, "-", 4)
	if len(rolls) != 2 {
		t.Errorf("expected 2 roll results, but found %d results\n", len(rolls))
	}
	//check sum
	expectedSum := 0
	for _, roll := range rolls {
		expectedSum += roll
	}
	if expectedSum != sum {
		t.Errorf("expected roll results to equal sum when added, but roll results equaled %d and sum was %d\n", expectedSum, sum)
	}
	//check modified sum, it should be sum minus four
	if modSum != (sum - 4) {
		t.Errorf("expected modified sum to be sum - 4, but modified sum equaled %d and sum was %d\n", modSum, sum)
	}
}

func TestRollWithUnsupportedModifier(t *testing.T) {
	//roll two ten sided dice and pass an unsupported operator
	rolls, sum, modSum := RollAndModify(2, 10, "z", 4)
	if len(rolls) != 2 {
		t.Errorf("expected 2 roll results, but found %d results\n", len(rolls))
	}
	//check sum
	expectedSum := 0
	for _, roll := range rolls {
		expectedSum += roll
	}
	if expectedSum != sum {
		t.Errorf("expected roll results to equal sum when added, but roll results equaled %d and sum was %d\n", expectedSum, sum)
	}
	//check modified sum, it should equal sum because operator is unsupported
	if modSum != sum {
		t.Errorf("expected modified sum to equal sum, but modified sum equaled %d and sum was %d\n", modSum, sum)
	}
}

func TestModifyDirectly(t *testing.T) {
	modifiedValue := Modify(2, "+", 2)
	expected := 4
	if modifiedValue != expected {
		t.Errorf("+ expected %d, actual %d", expected, modifiedValue)
	}
	modifiedValue = Modify(8, "-", 5)
	expected = 3
	if modifiedValue != expected {
		t.Errorf("- expected %d, actual %d", expected, modifiedValue)
	}
	modifiedValue = Modify(6, "?", 3)
	expected = 6
	if modifiedValue != expected {
		t.Errorf("? expected %d, actual %d", expected, modifiedValue)
	}
}

func TestRollMin(t *testing.T) {
	//roll 7d20 and provide min (lowest) roll
	rolls, min := RollMin(7, 20)

	//seven rolls should be present
	if len(rolls) != 7 {
		t.Errorf("expected 7 roll results, but found %d results\n", len(rolls))
	}
	//check min
	expectedMin := rolls[0]
	for _, roll := range rolls {
		if roll < expectedMin {
			expectedMin = roll
		}
	}
	if expectedMin != min {
		t.Errorf("expected min to be %d, but min was %d\n", expectedMin, min)
	}
	//rolls should be 1, 2, 3, 4, 5, ... 20
	for _, roll := range rolls {
		if roll < 1 || roll > 20 {
			t.Errorf("expected roll result to be 1, 2, 3, 4, 5, ... 20, but roll was %d\n", roll)
		}
	}
}

func TestRollMax(t *testing.T) {
	//roll 7d20 and provide max (highest) roll
	rolls, max := RollMax(7, 20)

	//seven rolls should be present
	if len(rolls) != 7 {
		t.Errorf("expected 7 roll results, but found %d results\n", len(rolls))
	}
	//check max
	expectedMax := rolls[0]
	for _, roll := range rolls {
		if roll > expectedMax {
			expectedMax = roll
		}
	}
	if expectedMax != max {
		t.Errorf("expected max to be %d, but max was %d\n", expectedMax, max)
	}
	//rolls should be 1, 2, 3, 4, 5, ... 20
	for _, roll := range rolls {
		if roll < 1 || roll > 20 {
			t.Errorf("expected roll result to be 1, 2, 3, 4, 5, ... 20, but roll was %d\n", roll)
		}
	}
}

func TestRollRange(t *testing.T) {
	//roll 1000 four sided dice
	rolls, _ := Roll(1000, 4)
	//1000 rolls should be present
	if len(rolls) != 1000 {
		t.Errorf("expected 1000 roll results, but found %d results\n", len(rolls))
	}
	//rolls should be 1, 2, 3, or 4
	for _, roll := range rolls {
		if roll < 1 || roll > 4 {
			t.Errorf("expected roll result to be 1, 2, 3, or 4, but roll was %d\n", roll)
		}
	}
}

func Test_ValidRollExpression(t *testing.T) {
	testCases := []struct {
		expression string
		want       bool
	}{
		{
			expression: "3d4+8",
			want:       true,
		},
		{
			expression: "d8-1",
			want:       true,
		},
		{
			expression: "1d12+3+1d8",
			want:       true,
		},
		{
			expression: "d12+3-d8",
			want:       true,
		},
		{
			expression: "2d20+", //+0
			want:       true,
		},
		{
			expression: "2d20-", //-0
			want:       true,
		},
		{
			expression: "2d20+1+",
			want:       false,
		},
		{
			expression: "2d20+1-",
			want:       false,
		},
		{
			expression: "1Z12+3/1d8",
			want:       false,
		},
		{
			expression: "foo+notbar",
			want:       false,
		},
		{
			expression: "heyo",
			want:       false,
		},
		{
			expression: "Roll a 2d6 and 3d12+3",
			want:       false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d) %s", i, tc.expression), func(t *testing.T) {
			got := ValidRollExpression(tc.expression)
			if got != tc.want {
				t.Errorf("want %t, got %t", tc.want, got)
			}
		})
	}
}

func Test_ContainsValidRollExpression(t *testing.T) {
	testCases := []struct {
		text string
		want int
	}{
		{
			text: "3d4+8",
			want: 1,
		},
		{
			text: "2W20+3",
			want: 0,
		},
		{
			text: "foobar",
			want: 0,
		},
		{
			text: "Roll a 2d6 and 3d12+3",
			want: 2,
		},
		{
			text: "These go together 2d12+3+d8, but not with this {{2d6}}",
			want: 2,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d) %s", i, tc.text), func(t *testing.T) {
			got := ContainsValidRollExpression(tc.text)
			if got != tc.want {
				t.Errorf("want %d, got %d", tc.want, got)
			}
		})
	}
}

func Test_RollExpression(t *testing.T) {
	simpleCases := []struct {
		expression string
		rolls      []int
		sum        int
		err        error
	}{
		{
			expression: "2d1+3",
			rolls:      []int{1, 1},
			sum:        5,
			err:        nil,
		},
		{
			expression: "3d1",
			rolls:      []int{1, 1, 1},
			sum:        3,
			err:        nil,
		},
		{
			expression: "3d1+2d1+3",
			rolls:      []int{1, 1, 1, 1, 1},
			sum:        8,
			err:        nil,
		},
		{
			expression: "heyo",
			rolls:      nil,
			sum:        0,
			err:        ErrInvalidRollExpression,
		},
		{
			expression: "a5d10*zz",
			rolls:      nil,
			sum:        0,
			err:        ErrInvalidRollExpression,
		},
	}

	for i, tc := range simpleCases {
		t.Run(fmt.Sprintf("%d) %s", i, tc.expression), func(t *testing.T) {
			rolls, sum, err := RollExpression(tc.expression)
			if !reflect.DeepEqual(rolls, tc.rolls) {
				t.Errorf("want %v, got %v", tc.rolls, rolls)
			}
			if sum != tc.sum {
				t.Errorf("want %d, got %d", tc.sum, sum)
			}
			if err != tc.err {
				t.Errorf("want %s, got %s", tc.err, err)
			}
		})
	}

	detailedCases := []struct {
		expression  string //expression to test
		subtractDie int    //if expression pair and second is to be subtracted number of die in second expression
		modifer     int    //the total value of the modifiers, negative numbers work also
		rollLen     int    //the total number of die rolled
		rollMin     int    //the lowest possible die number
		rollMax     int    //the highest possible die number
		err         error  //error if one is expected
	}{
		{
			expression: "5d10",
			rollLen:    5,
			rollMin:    1,
			rollMax:    10,
			err:        nil,
		},
		{
			expression: "d10",
			rollLen:    1,
			rollMin:    1,
			rollMax:    10,
			err:        nil,
		},
		{
			expression: "1d20+6",
			modifer:    6,
			rollLen:    1,
			rollMin:    1,
			rollMax:    20,
			err:        nil,
		},
		{
			expression: "3d33+",
			modifer:    0,
			rollLen:    3,
			rollMin:    1,
			rollMax:    33,
			err:        nil,
		},
		{
			expression: "6d12-",
			modifer:    0,
			rollLen:    6,
			rollMin:    1,
			rollMax:    12,
			err:        nil,
		},
		{
			expression: "2d20+3",
			modifer:    3,
			rollLen:    2,
			rollMin:    1,
			rollMax:    20,
			err:        nil,
		},
		{
			expression: "3d6-2",
			modifer:    -2,
			rollLen:    3,
			rollMin:    1,
			rollMax:    6,
			err:        nil,
		},
		{
			expression: "2d1+3",
			modifer:    3,
			rollLen:    2,
			rollMin:    1,
			rollMax:    1,
			err:        nil,
		},
		{
			expression: "2d20+3+2d6+1",
			modifer:    4,
			rollLen:    4,
			rollMin:    1,
			rollMax:    20,
			err:        nil,
		},
		{
			expression:  "d12+3-2d4",
			subtractDie: 2,
			modifer:     3,
			rollLen:     3,
			rollMin:     1,
			rollMax:     12,
			err:         nil,
		},
		{
			expression: "d12+d8",
			rollLen:    2,
			rollMin:    1,
			rollMax:    12,
			err:        nil,
		},
		{
			expression:  "d12-d3",
			subtractDie: 1,
			rollLen:     2,
			rollMin:     1,
			rollMax:     12,
			err:         nil,
		},
		{
			expression: "min:2d20+1d6",
			rollLen:    0,
			err:        ErrInvalidRollExpression,
		},
		{
			expression: "heyo",
			rollLen:    0,
			err:        ErrInvalidRollExpression,
		},
	}

	for i, tc := range detailedCases {
		t.Run(fmt.Sprintf("%d) %s", i, tc.expression), func(t *testing.T) {
			rolls, sum, err := RollExpression(tc.expression)
			if len(rolls) != tc.rollLen {
				t.Errorf("[len] want %d, got %d", tc.rollLen, len(rolls))
			}

			wantSum := 0
			for r, roll := range rolls {
				if tc.subtractDie > 0 && r >= (tc.rollLen-tc.subtractDie) {
					wantSum -= roll
					continue
				}
				wantSum += roll
			}
			wantSum += tc.modifer
			if wantSum != sum {
				t.Errorf("[sum] want %d, got %d", wantSum, sum)
			}

			for _, roll := range rolls {
				if roll < tc.rollMin || roll > tc.rollMax {
					t.Errorf("[rolls] want roll %d-%d, got %d", tc.rollMin, tc.rollMax, roll)
				}
			}

			if err != tc.err {
				t.Errorf("[err] want %s, got %s", tc.err, err)
			}
		})
	}
}

func TestRollInvalidExpression(t *testing.T) {
	rolls, sum, err := RollExpression("a5d10*zz")
	if err == nil {
		t.Errorf("error was expected, but didn't receive an error\n")
	}
	if rolls != nil {
		t.Errorf("rolls should be nil, but received %v\n", rolls)
	}
	if sum != 0 {
		t.Errorf("expected sum to be 0, but sum was %d\n", sum)
	}
}

func TestRollExpressionMinPrefix(t *testing.T) {
	rolls, min, err := RollExpression("min:2d20")
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}
	//two rolls should be present
	if len(rolls) != 2 {
		t.Errorf("expected 2 roll results, but found %d results\n", len(rolls))
	}
	//check min
	expectedMin := rolls[0]
	for _, roll := range rolls {
		if roll < expectedMin {
			expectedMin = roll
		}
	}
	if expectedMin != min {
		t.Errorf("expected min to be %d, but min was %d\n", expectedMin, min)
	}
	//rolls should be 1, 2, 3, 4, 5, ... 20
	for _, roll := range rolls {
		if roll < 1 || roll > 20 {
			t.Errorf("expected roll result to be 1, 2, 3, 4, 5, ... 20, but roll was %d\n", roll)
		}
	}
}

func TestRollExpressionMinPrefixWithModifier(t *testing.T) {
	rolls, min, err := RollExpression("min:2d20+3")
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}
	//two rolls should be present
	if len(rolls) != 2 {
		t.Errorf("expected 2 roll results, but found %d results\n", len(rolls))
	}
	//check min
	expectedMin := rolls[0]
	for _, roll := range rolls {
		if roll < expectedMin {
			expectedMin = roll
		}
	}
	expectedMin += 3
	if expectedMin != min {
		t.Errorf("expected min to be %d, but min was %d\n", expectedMin, min)
	}
	//rolls should be 1, 2, 3, 4, 5, ... 20
	for _, roll := range rolls {
		if roll < 1 || roll > 20 {
			t.Errorf("expected roll result to be 1, 2, 3, 4, 5, ... 20, but roll was %d\n", roll)
		}
	}
}

func TestRollExpressionMaxPrefix(t *testing.T) {
	rolls, max, err := RollExpression("max:2d20")
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}
	//two rolls should be present
	if len(rolls) != 2 {
		t.Errorf("expected 2 roll results, but found %d results\n", len(rolls))
	}
	//check max
	expectedMax := rolls[0]
	for _, roll := range rolls {
		if roll > expectedMax {
			expectedMax = roll
		}
	}
	if expectedMax != max {
		t.Errorf("expected max to be %d, but max was %d\n", expectedMax, max)
	}
	//rolls should be 1, 2, 3, 4, 5, ... 20
	for _, roll := range rolls {
		if roll < 1 || roll > 20 {
			t.Errorf("expected roll result to be 1, 2, 3, 4, 5, ... 20, but roll was %d\n", roll)
		}
	}
}

func TestRollExpressionMaxPrefixWithModifier(t *testing.T) {
	rolls, max, err := RollExpression("max:2d20-1")
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}
	//two rolls should be present
	if len(rolls) != 2 {
		t.Errorf("expected 2 roll results, but found %d results\n", len(rolls))
	}
	//check max
	expectedMax := rolls[0]
	for _, roll := range rolls {
		if roll > expectedMax {
			expectedMax = roll
		}
	}
	expectedMax--
	if expectedMax != max {
		t.Errorf("expected max to be %d, but max was %d\n", expectedMax, max)
	}
	//rolls should be 1, 2, 3, 4, 5, ... 20
	for _, roll := range rolls {
		if roll < 1 || roll > 20 {
			t.Errorf("expected roll result to be 1, 2, 3, 4, 5, ... 20, but roll was %d\n", roll)
		}
	}
}

func TestRollExpressionDubPrefix(t *testing.T) {
	rolls, sum, err := RollExpression("dub:3d8")
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}
	//three rolls should be present
	if len(rolls) != 3 {
		t.Errorf("expected 3 roll results, but found %d results\n", len(rolls))
	}
	//check sum
	expectedSum := 0
	for _, roll := range rolls {
		expectedSum += roll
	}
	expectedSum *= 2
	if expectedSum != sum {
		t.Errorf("expected sum to be %d, but sum was %d\n", expectedSum, sum)
	}
	//rolls should be 1, 2, 3, 4, 5, ... 8
	for _, roll := range rolls {
		if roll < 1 || roll > 8 {
			t.Errorf("expected roll result to be 1, 2, 3, 4, 5, ... 8, but roll was %d\n", roll)
		}
	}
}

func TestRollExpressionHalfPrefix(t *testing.T) {
	rolls, sum, err := RollExpression("half:4d13")
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}
	//four rolls should be present
	if len(rolls) != 4 {
		t.Errorf("expected 4 roll results, but found %d results\n", len(rolls))
	}
	//check sum
	expectedSum := 0
	for _, roll := range rolls {
		expectedSum += roll
	}
	expectedSum /= 2
	if expectedSum != sum {
		t.Errorf("expected sum to be %d, but sum was %d\n", expectedSum, sum)
	}
	//rolls should be 1, 2, 3, 4, 5, ... 13
	for _, roll := range rolls {
		if roll < 1 || roll > 13 {
			t.Errorf("expected roll result to be 1, 2, 3, 4, 5, ... 13, but roll was %d\n", roll)
		}
	}
}

func TestRollExpressionDropLowestPrefix(t *testing.T) {
	rolls, sum, err := RollExpression("dropL:4d6+2")
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}
	//four rolls should be present
	if len(rolls) != 4 {
		t.Errorf("expected 4 roll results, but found %d results\n", len(rolls))
	}
	//check sum
	expectedSum := 0
	lowest := rolls[0]
	for _, roll := range rolls {
		expectedSum += roll
		lowest = min(lowest, roll)
	}
	expectedSum += 2
	expectedSum -= lowest
	if expectedSum != sum {
		t.Errorf("expected sum to be %d, but sum was %d\n", expectedSum, sum)
	}
	//rolls should be 1, 2, 3, 4, 5, or 6
	for _, roll := range rolls {
		if roll < 1 || roll > 6 {
			t.Errorf("expected roll result to be 1, 2, 3, 4, 5, or 6, but roll was %d\n", roll)
		}
	}
}

func TestRollExpressionDropHighestPrefix(t *testing.T) {
	rolls, sum, err := RollExpression("dropH:4d6+2")
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}
	//four rolls should be present
	if len(rolls) != 4 {
		t.Errorf("expected 4 roll results, but found %d results\n", len(rolls))
	}
	//check sum
	expectedSum := 0
	highest := rolls[0]
	for _, roll := range rolls {
		expectedSum += roll
		highest = max(highest, roll)
	}
	expectedSum += 2
	expectedSum -= highest
	if expectedSum != sum {
		t.Errorf("expected sum to be %d, but sum was %d\n", expectedSum, sum)
	}
	//rolls should be 1, 2, 3, 4, 5, or 6
	for _, roll := range rolls {
		if roll < 1 || roll > 6 {
			t.Errorf("expected roll result to be 1, 2, 3, 4, 5, or 6, but roll was %d\n", roll)
		}
	}
}

func TestRollSimplestExpression(t *testing.T) {
	rolls, sum, err := RollExpression("d6")
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}
	//one roll should be present
	if len(rolls) != 1 {
		t.Errorf("expected 1 roll result, but found %d results\n", len(rolls))
	}
	//single roll means sum and roll should be equal
	if rolls[0] != sum {
		t.Errorf("expected sum to be %d, but sum was %d\n", rolls[0], sum)
	}
	//rolls should be 1, 2, 3, 4, 5, or 6
	for _, roll := range rolls {
		if roll < 1 || roll > 6 {
			t.Errorf("expected roll result to be 1, 2, 3, 4, 5, or 6, but roll was %d\n", roll)
		}
	}
}

func TestAFewRollExpressions(t *testing.T) {
	rolls, sum, err := RollExpression("d6")
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}
	//one roll should be present
	if len(rolls) != 1 {
		t.Errorf("expected 1 roll result, but found %d results\n", len(rolls))
	}
	//single roll means sum and roll should be equal
	if rolls[0] != sum {
		t.Errorf("expected sum to be %d, but sum was %d\n", rolls[0], sum)
	}
	//rolls should be 1, 2, 3, 4, 5, or 6
	for _, roll := range rolls {
		if roll < 1 || roll > 6 {
			t.Errorf("expected roll result to be 1, 2, 3, 4, 5, or 6, but roll was %d\n", roll)
		}
	}

	rolls, sum, err = RollExpression("7d16")
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}
	//seven roll should be present
	if len(rolls) != 7 {
		t.Errorf("expected 7 roll results, but found %d results\n", len(rolls))
	}
	//check sum
	expectedSum := 0
	for _, roll := range rolls {
		expectedSum += roll
	}
	if expectedSum != sum {
		t.Errorf("expected sum to be %d, but sum was %d\n", expectedSum, sum)
	}
	//rolls should be 1, 2, 3, 4, 5, ... 16
	for _, roll := range rolls {
		if roll < 1 || roll > 16 {
			t.Errorf("expected roll result to be 1, 2, 3, 4, 5, ... 16, but roll was %d\n", roll)
		}
	}

	rolls, sum, err = RollExpression("3d20-4")
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}
	//three roll should be present
	if len(rolls) != 3 {
		t.Errorf("expected 3 roll results, but found %d results\n", len(rolls))
	}
	//check sum
	expectedSum = 0
	for _, roll := range rolls {
		expectedSum += roll
	}
	expectedSum -= 4
	if expectedSum != sum {
		t.Errorf("expected sum to be %d, but sum was %d\n", expectedSum, sum)
	}
	//rolls should be 1, 2, 3, 4, 5, ... 20
	for _, roll := range rolls {
		if roll < 1 || roll > 20 {
			t.Errorf("expected roll result to be 1, 2, 3, 4, 5, ... 20, but roll was %d\n", roll)
		}
	}
}

func TestRollChallengeMustBeGreater(t *testing.T) {
	succeeded, result, found, err := RollChallenge("1d20+3", 13, false, nil)
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}

	if succeeded && (result <= 13) {
		t.Errorf("expected to fail challenge because %d is not greater than 13, but succeeded\n", result)
	}

	if found != nil {
		t.Errorf("found should be nil, but received %v\n", found)
	}
}
func TestRollChallengeCanBeEqual(t *testing.T) {
	succeeded, result, found, err := RollChallenge("1d20+3", 13, true, nil)
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}

	if succeeded && (result < 13) {
		t.Errorf("expected to fail challenge because %d is less than 13, but succeeded\n", result)
	}

	if found != nil {
		t.Errorf("found should be nil, but received %v\n", found)
	}
}

func TestRollChallengeWhenEqualsAndMustBeGreater(t *testing.T) {
	succeeded, result, _, err := RollChallenge("1d1+3", 4, false, nil)
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}

	if succeeded {
		t.Errorf("expected to fail challenge because %d is not greater than 4, but succeeded\n", result)
	}
}

func TestRollChallengeWhenEqualsAndCanBeEqual(t *testing.T) {
	succeeded, result, _, err := RollChallenge("1d1+3", 4, true, nil)
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}

	if !succeeded {
		t.Errorf("expected to succeed challenge because %d is equal to 4, but failed\n", result)
	}
}

func TestRollChallengeWithAlerts(t *testing.T) {
	//this should always result in a roll of 1 or 2 which means found will have a value because we are alerting on 1 or 2
	//the modifier allows for a chance to fail the challenge
	succeeded, result, found, err := RollChallenge("1d2+3", 5, true, []int{1, 2})
	if err != nil {
		t.Errorf("error was not expected, but err was encountered %s\n", err)
	}

	if succeeded && (result < 5) {
		t.Errorf("expected to fail challenge because %d is less than 5, but succeeded\n", result)
	}

	if len(found) != 1 {
		t.Errorf("expected to find a single value in found, but it contained %d\n", len(found))
	}

	if found[0] != (result - 3) {
		t.Errorf("expected found to be 1 or 2, but found was %d\n", found[0])
	}
}

func TestRollChallengeWithInvalidExpression(t *testing.T) {
	_, _, _, err := RollChallenge("foobar", 13, false, nil)
	if err == nil {
		t.Errorf("error was expected, but didn't receive an error\n")
	}
}

func Test_RollString(t *testing.T) {
	testCases := []struct {
		name    string
		rollStr string
		want    string
	}{
		{
			name:    "validate roll is replaced with valid value...",
			rollStr: "This should be {{1d1+3}}.",
			want:    "This should be 4.",
		},
		{
			name:    "validate roll is replaced with valid values...",
			rollStr: "This should be {{1d1+3}} and {{2d1}}. Right?",
			want:    "This should be 4 and 2. Right?",
		},
		{
			name:    "validate value is unchanged if no roll expression in string...",
			rollStr: "This should be the same!",
			want:    "This should be the same!",
		},
		{
			name:    "validate value is unchanged roll expression invalid...",
			rollStr: "This should be {{1dbroke}} the same!",
			want:    "This should be {{1dbroke}} the same!",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got := RollString(test.rollStr)

			if got != test.want {
				t.Errorf("want %s, got %s", test.want, got)
			}
		})
	}
}
