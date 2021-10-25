package dice

import (
	"fmt"
	"testing"
)

func Test_Roll(t *testing.T) {
	testCases := []struct {
		number  int
		sides   int
		rollLen int
		rollMin int
		rollMax int
	}{
		{
			number:  1,
			sides:   6,
			rollLen: 1,
			rollMin: 1,
			rollMax: 6,
		},
		{
			number:  4,
			sides:   6,
			rollLen: 4,
			rollMin: 1,
			rollMax: 6,
		},
		{
			number:  3,
			sides:   20,
			rollLen: 3,
			rollMin: 1,
			rollMax: 20,
		},
		{
			number:  1000,
			sides:   4,
			rollLen: 1000,
			rollMin: 1,
			rollMax: 4,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d) number %d, sides %d", i, tc.number, tc.sides), func(t *testing.T) {
			rolls, sum := Roll(tc.number, tc.sides)
			if len(rolls) != tc.rollLen {
				t.Errorf("[len] want %d, got %d", tc.rollLen, len(rolls))
			}

			wantSum := 0
			for _, roll := range rolls {
				wantSum += roll
			}
			if wantSum != sum {
				t.Errorf("[sum] want %d, got %d", wantSum, sum)
			}

			for _, roll := range rolls {
				if roll < tc.rollMin || roll > tc.rollMax {
					t.Errorf("[rolls] want roll %d-%d, got %d", tc.rollMin, tc.rollMax, roll)
				}
			}
		})
	}
}

func Test_RollAndModify(t *testing.T) {
	testCases := []struct {
		number   int
		sides    int
		operator string
		modifer  int
		rollLen  int
		rollMin  int
		rollMax  int
		err      error
	}{
		{
			number:   2,
			sides:    6,
			operator: "+",
			modifer:  3,
			rollLen:  2,
			rollMin:  1,
			rollMax:  6,
			err:      nil,
		},
		{
			number:   2,
			sides:    10,
			operator: "-",
			modifer:  4,
			rollLen:  2,
			rollMin:  1,
			rollMax:  10,
			err:      nil,
		},
		{
			number:   2,
			sides:    20,
			operator: "Z",
			modifer:  3,
			rollLen:  0,
			rollMin:  0,
			rollMax:  0,
			err:      ErrInvalidOperator,
		},
		{
			number:   3,
			sides:    4,
			operator: "-",
			modifer:  12,
			rollLen:  3,
			rollMin:  1,
			rollMax:  4,
			err:      nil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d) number %d, sides %d, operator %s, modifier %d", i, tc.number, tc.sides, tc.operator, tc.modifer), func(t *testing.T) {
			rolls, sum, modSum, err := RollAndModify(tc.number, tc.sides, tc.operator, tc.modifer)
			if len(rolls) != tc.rollLen {
				t.Errorf("[len] want %d, got %d", tc.rollLen, len(rolls))
			}

			wantSum := 0
			for _, roll := range rolls {
				wantSum += roll
			}
			if wantSum != sum {
				t.Errorf("[sum] want %d, got %d", wantSum, sum)
			}

			wantModSum := wantSum
			switch tc.operator {
			case "+":
				wantModSum += tc.modifer
			case "-":
				wantModSum -= tc.modifer
			}
			if wantModSum != modSum {
				t.Errorf("[mod sum] want %d, got %d", wantModSum, modSum)
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

func Test_Modify(t *testing.T) {
	testCases := []struct {
		value    int
		operator string
		modifer  int
		want     int
		err      error
	}{
		{
			value:    2,
			operator: "+",
			modifer:  2,
			want:     4,
			err:      nil,
		},
		{
			value:    8,
			operator: "-",
			modifer:  5,
			want:     3,
			err:      nil,
		},
		{
			value:    6,
			operator: "?",
			modifer:  1,
			want:     0,
			err:      ErrInvalidOperator,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d) value %d, operator %s, modifer %d", i, tc.value, tc.operator, tc.modifer), func(t *testing.T) {
			got, err := Modify(tc.value, tc.operator, tc.modifer)

			if got != tc.want {
				t.Errorf("[mod sum] want %d, got %d", tc.want, got)
			}

			if err != tc.err {
				t.Errorf("[err] want %s, got %s", tc.err, err)
			}
		})
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
			expression: "2d20+", //acts as +0
			want:       true,
		},
		{
			expression: "2d20-", //acts as -0
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
	detailedCases := []struct {
		expression               string //expression to test
		secondExpressionDieCount int    //the number of die in the second expression if present
		subtractDie              bool   //should the second expression die be subtracted
		modifer                  int    //the total value of the modifiers, negative numbers work also
		rollLen                  int    //the total number of die rolled
		rollMin                  int    //the lowest possible die number
		rollMax                  int    //the highest possible die number
		err                      error  //error if one is expected
	}{
		{
			expression: "d6",
			rollLen:    1,
			rollMin:    1,
			rollMax:    6,
			err:        nil,
		},
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
			expression:               "d12+3-2d4",
			secondExpressionDieCount: 2,
			subtractDie:              true,
			modifer:                  3,
			rollLen:                  3,
			rollMin:                  1,
			rollMax:                  12,
			err:                      nil,
		},
		{
			expression: "d12+d8",
			rollLen:    2,
			rollMin:    1,
			rollMax:    12,
			err:        nil,
		},
		{
			expression:               "d12-d3",
			secondExpressionDieCount: 1,
			subtractDie:              true,
			rollLen:                  2,
			rollMin:                  1,
			rollMax:                  12,
			err:                      nil,
		},
		{
			expression: "heyo",
			rollLen:    0,
			err:        ErrInvalidRollExpression,
		},
		{
			expression: "a5d10*zz",
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
				if tc.subtractDie && r >= (tc.rollLen-tc.secondExpressionDieCount) {
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

	//prefix tests
	maxTestCases := []struct {
		expression string
		modifer    int
		rollLen    int
		rollMin    int
		rollMax    int
		err        error
	}{
		{
			expression: "max:2d20",
			rollLen:    2,
			rollMin:    1,
			rollMax:    20,
			err:        nil,
		},
		{
			expression: "max:2d20-2",
			modifer:    -2,
			rollLen:    2,
			rollMin:    1,
			rollMax:    20,
			err:        nil,
		},
		{
			expression: "max:2d20+1d6",
			rollLen:    0,
			rollMin:    0,
			rollMax:    0,
			err:        ErrInvalidRollExpression,
		},
		{
			expression: "max:2d20zzz",
			rollLen:    0,
			rollMin:    0,
			rollMax:    0,
			err:        ErrInvalidRollExpression,
		},
	}

	for i, tc := range maxTestCases {
		t.Run(fmt.Sprintf("%d) %s", i, tc.expression), func(t *testing.T) {
			rolls, max, err := RollExpression(tc.expression)
			if len(rolls) != tc.rollLen {
				t.Errorf("[len] want %d, got %d", tc.rollLen, len(rolls))
			}

			wantMax := 0
			for i, roll := range rolls {
				if i == 0 {
					wantMax = roll
					continue
				}
				if roll > wantMax {
					wantMax = roll
				}
			}
			wantMax += tc.modifer
			if wantMax != max {
				t.Errorf("[max] want %d, got %d", wantMax, max)
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

	minTestCases := []struct {
		expression string
		modifer    int
		rollLen    int
		rollMin    int
		rollMax    int
		err        error
	}{
		{
			expression: "min:2d20",
			rollLen:    2,
			rollMin:    1,
			rollMax:    20,
			err:        nil,
		},
		{
			expression: "min:2d20+3",
			modifer:    3,
			rollLen:    2,
			rollMin:    1,
			rollMax:    20,
			err:        nil,
		},
		{
			expression: "min:2d20+1d6",
			rollLen:    0,
			rollMin:    0,
			rollMax:    0,
			err:        ErrInvalidRollExpression,
		},
		{
			expression: "min:2E20fff",
			rollLen:    0,
			rollMin:    0,
			rollMax:    0,
			err:        ErrInvalidRollExpression,
		},
	}

	for i, tc := range minTestCases {
		t.Run(fmt.Sprintf("%d) %s", i, tc.expression), func(t *testing.T) {
			rolls, min, err := RollExpression(tc.expression)
			if len(rolls) != tc.rollLen {
				t.Errorf("[len] want %d, got %d", tc.rollLen, len(rolls))
			}

			wantMin := 0
			for i, roll := range rolls {
				if i == 0 {
					wantMin = roll
					continue
				}
				if roll < wantMin {
					wantMin = roll
				}
			}
			wantMin += tc.modifer
			if wantMin != min {
				t.Errorf("[min] want %d, got %d", wantMin, min)
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

	dubTestCases := []struct {
		expression               string
		secondExpressionDieCount int
		subtractDie              bool
		modifer                  int
		rollLen                  int
		rollMin                  int
		rollMax                  int
		err                      error
	}{
		{
			expression: "dub:3d8",
			rollLen:    3,
			rollMin:    1,
			rollMax:    8,
			err:        nil,
		},
		{
			expression: "dub:3d8+3",
			modifer:    3,
			rollLen:    3,
			rollMin:    1,
			rollMax:    8,
			err:        nil,
		},
		{
			expression: "dub:3d8+3+d6",
			modifer:    3,
			rollLen:    4,
			rollMin:    1,
			rollMax:    8,
			err:        nil,
		},
		{
			expression:               "dub:3d20+3-2d6",
			secondExpressionDieCount: 2,
			subtractDie:              true,
			modifer:                  3,
			rollLen:                  5,
			rollMin:                  1,
			rollMax:                  20,
			err:                      nil,
		},
		{
			expression: "dub:2E20fff",
			rollLen:    0,
			rollMin:    0,
			rollMax:    0,
			err:        ErrInvalidRollExpression,
		},
	}

	for i, tc := range dubTestCases {
		t.Run(fmt.Sprintf("%d) %s", i, tc.expression), func(t *testing.T) {
			rolls, sum, err := RollExpression(tc.expression)
			if len(rolls) != tc.rollLen {
				t.Errorf("[len] want %d, got %d", tc.rollLen, len(rolls))
			}

			wantSum := 0
			for r, roll := range rolls {
				if tc.subtractDie && r >= (tc.rollLen-tc.secondExpressionDieCount) {
					wantSum -= roll
					continue
				}
				wantSum += roll
			}
			wantSum += tc.modifer
			wantSum *= 2
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

	halfTestCases := []struct {
		expression               string
		secondExpressionDieCount int
		subtractDie              bool
		modifer                  int
		rollLen                  int
		rollMin                  int
		rollMax                  int
		err                      error
	}{
		{
			expression: "half:4d13",
			rollLen:    4,
			rollMin:    1,
			rollMax:    13,
			err:        nil,
		},
		{
			expression: "half:3d8+3",
			modifer:    3,
			rollLen:    3,
			rollMin:    1,
			rollMax:    8,
			err:        nil,
		},
		{
			expression: "half:3d8+3+d6",
			modifer:    3,
			rollLen:    4,
			rollMin:    1,
			rollMax:    8,
			err:        nil,
		},
		{
			expression:               "half:3d20+3-2d6",
			secondExpressionDieCount: 2,
			subtractDie:              true,
			modifer:                  3,
			rollLen:                  5,
			rollMin:                  1,
			rollMax:                  20,
			err:                      nil,
		},
		{
			expression: "half:2E20fff",
			rollLen:    0,
			rollMin:    0,
			rollMax:    0,
			err:        ErrInvalidRollExpression,
		},
	}

	for i, tc := range halfTestCases {
		t.Run(fmt.Sprintf("%d) %s", i, tc.expression), func(t *testing.T) {
			rolls, sum, err := RollExpression(tc.expression)
			if len(rolls) != tc.rollLen {
				t.Errorf("[len] want %d, got %d", tc.rollLen, len(rolls))
			}

			wantSum := 0
			for r, roll := range rolls {
				if tc.subtractDie && r >= (tc.rollLen-tc.secondExpressionDieCount) {
					wantSum -= roll
					continue
				}
				wantSum += roll
			}
			wantSum += tc.modifer
			wantSum /= 2
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

	dropLTestCases := []struct {
		expression               string
		secondExpressionDieCount int
		subtractDie              bool
		modifer                  int
		secondModifier           int
		rollLen                  int
		rollMin                  int
		rollMax                  int
		err                      error
	}{
		{
			expression: "dropL:4d13",
			rollLen:    4,
			rollMin:    1,
			rollMax:    13,
			err:        nil,
		},
		{
			expression: "dropL:3d8+3",
			modifer:    3,
			rollLen:    3,
			rollMin:    1,
			rollMax:    8,
			err:        nil,
		},
		{
			expression:               "dropL:3d8+3+d6",
			secondExpressionDieCount: 1,
			modifer:                  3,
			rollLen:                  4,
			rollMin:                  1,
			rollMax:                  8,
			err:                      nil,
		},
		{
			expression:               "dropL:2d12+4-d6-2",
			secondExpressionDieCount: 1,
			subtractDie:              true,
			modifer:                  4,
			secondModifier:           -2,
			rollLen:                  3,
			rollMin:                  1,
			rollMax:                  12,
			err:                      nil,
		},
		{
			expression:               "dropL:3d20+3-2d6",
			secondExpressionDieCount: 2,
			subtractDie:              true,
			modifer:                  3,
			rollLen:                  5,
			rollMin:                  1,
			rollMax:                  20,
			err:                      nil,
		},
		{
			expression: "dropL:2E20fff",
			rollLen:    0,
			rollMin:    0,
			rollMax:    0,
			err:        ErrInvalidRollExpression,
		},
	}

	for i, tc := range dropLTestCases {
		t.Run(fmt.Sprintf("%d) %s", i, tc.expression), func(t *testing.T) {
			rolls, sum, err := RollExpression(tc.expression)
			if len(rolls) != tc.rollLen {
				t.Errorf("[len] want %d, got %d", tc.rollLen, len(rolls))
			}

			//handle sum for main expression
			wantSum := 0
			lowest := 0
			for r, roll := range rolls[:(tc.rollLen - tc.secondExpressionDieCount)] {
				if r == 0 {
					lowest = roll
				} else if roll < lowest {
					lowest = roll
				}

				wantSum += roll
			}

			//handle second expression if present
			wantSumSecond := 0
			for _, roll := range rolls[(tc.rollLen - tc.secondExpressionDieCount):] {
				wantSumSecond += roll
			}
			wantSumSecond += tc.secondModifier

			wantSum += tc.modifer
			wantSum -= lowest

			if tc.subtractDie {
				wantSum -= wantSumSecond
			} else {
				wantSum += wantSumSecond
			}

			if wantSum != sum {
				t.Errorf("[sum rolls] %v", rolls)
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

	dropHTestCases := []struct {
		expression               string
		secondExpressionDieCount int
		subtractDie              bool
		modifer                  int
		secondModifier           int
		rollLen                  int
		rollMin                  int
		rollMax                  int
		err                      error
	}{
		{
			expression: "dropH:4d13",
			rollLen:    4,
			rollMin:    1,
			rollMax:    13,
			err:        nil,
		},
		{
			expression: "dropH:3d8+3",
			modifer:    3,
			rollLen:    3,
			rollMin:    1,
			rollMax:    8,
			err:        nil,
		},
		{
			expression:               "dropH:3d8+3+d6",
			secondExpressionDieCount: 1,
			modifer:                  3,
			rollLen:                  4,
			rollMin:                  1,
			rollMax:                  8,
			err:                      nil,
		},
		{
			expression:               "dropH:2d12+4-d6-2",
			secondExpressionDieCount: 1,
			subtractDie:              true,
			modifer:                  4,
			secondModifier:           -2,
			rollLen:                  3,
			rollMin:                  1,
			rollMax:                  12,
			err:                      nil,
		},
		{
			expression:               "dropH:3d20+3-2d6",
			secondExpressionDieCount: 2,
			subtractDie:              true,
			modifer:                  3,
			rollLen:                  5,
			rollMin:                  1,
			rollMax:                  20,
			err:                      nil,
		},
		{
			expression: "dropH:2E20fff",
			rollLen:    0,
			rollMin:    0,
			rollMax:    0,
			err:        ErrInvalidRollExpression,
		},
	}

	for i, tc := range dropHTestCases {
		t.Run(fmt.Sprintf("%d) %s", i, tc.expression), func(t *testing.T) {
			rolls, sum, err := RollExpression(tc.expression)
			if len(rolls) != tc.rollLen {
				t.Errorf("[len] want %d, got %d", tc.rollLen, len(rolls))
			}

			//handle sum for main expression
			wantSum := 0
			highest := 0
			for r, roll := range rolls[:(tc.rollLen - tc.secondExpressionDieCount)] {
				if r == 0 {
					highest = roll
				} else if roll > highest {
					highest = roll
				}

				wantSum += roll
			}

			//handle second expression if present
			wantSumSecond := 0
			for _, roll := range rolls[(tc.rollLen - tc.secondExpressionDieCount):] {
				wantSumSecond += roll
			}
			wantSumSecond += tc.secondModifier

			wantSum += tc.modifer
			wantSum -= highest

			if tc.subtractDie {
				wantSum -= wantSumSecond
			} else {
				wantSum += wantSumSecond
			}

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

func Test_RollMax(t *testing.T) {
	testCases := []struct {
		number  int
		sides   int
		rollLen int
		rollMin int
		rollMax int
	}{
		{
			number:  7,
			sides:   20,
			rollLen: 7,
			rollMin: 1,
			rollMax: 20,
		},
		{
			number:  3,
			sides:   100,
			rollLen: 3,
			rollMin: 1,
			rollMax: 100,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d) number %d, sides %d", i, tc.number, tc.sides), func(t *testing.T) {
			rolls, max := RollMax(tc.number, tc.sides)
			if len(rolls) != tc.rollLen {
				t.Errorf("[len] want %d, got %d", tc.rollLen, len(rolls))
			}

			wantMax := rolls[0]
			for _, roll := range rolls {
				if roll > wantMax {
					wantMax = roll
				}
			}
			if wantMax != max {
				t.Errorf("[max] want %d, got %d", wantMax, max)
			}

			for _, roll := range rolls {
				if roll < tc.rollMin || roll > tc.rollMax {
					t.Errorf("[rolls] want roll %d-%d, got %d", tc.rollMin, tc.rollMax, roll)
				}
			}
		})
	}
}
func Test_RollMin(t *testing.T) {
	testCases := []struct {
		number  int
		sides   int
		rollLen int
		rollMin int
		rollMax int
	}{
		{
			number:  7,
			sides:   20,
			rollLen: 7,
			rollMin: 1,
			rollMax: 20,
		},
		{
			number:  3,
			sides:   100,
			rollLen: 3,
			rollMin: 1,
			rollMax: 100,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d) number %d, sides %d", i, tc.number, tc.sides), func(t *testing.T) {
			rolls, min := RollMin(tc.number, tc.sides)
			if len(rolls) != tc.rollLen {
				t.Errorf("[len] want %d, got %d", tc.rollLen, len(rolls))
			}

			wantMin := rolls[0]
			for _, roll := range rolls {
				if roll < wantMin {
					wantMin = roll
				}
			}
			if wantMin != min {
				t.Errorf("[min] want %d, got %d", wantMin, min)
			}

			for _, roll := range rolls {
				if roll < tc.rollMin || roll > tc.rollMax {
					t.Errorf("[rolls] want roll %d-%d, got %d", tc.rollMin, tc.rollMax, roll)
				}
			}
		})
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
