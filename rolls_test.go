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
		err     error
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
		{
			number: 0,
			sides:  4,
		},
		{
			number: -1,
			sides:  6,
			err:    ErrInvalidNumberOfDice,
		},
		{
			number: 2,
			sides:  -6,
			err:    ErrInvalidNumberOfSides,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d) number %d, sides %d", i, tc.number, tc.sides), func(t *testing.T) {
			rolls, sum, err := Roll(tc.number, tc.sides)
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

			if err != tc.err {
				t.Errorf("[err] want %s, got %s", tc.err, err)
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
			err:      ErrInvalidOperator,
		},
		{
			number:   -1,
			sides:    4,
			operator: "+",
			err:      ErrInvalidNumberOfDice,
		},
		{
			number:   1,
			sides:    -4,
			operator: "-",
			err:      ErrInvalidNumberOfSides,
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
		{
			number:   0,
			sides:    4,
			operator: "+",
			modifer:  2,
			err:      nil,
		},
		{
			number:   0,
			sides:    4,
			operator: "-",
			modifer:  2,
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
		{
			value:    0,
			operator: "-",
			modifer:  1,
			want:     -1,
			err:      nil,
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

func Test_RollMax(t *testing.T) {
	testCases := []struct {
		number  int
		sides   int
		rollLen int
		rollMin int
		rollMax int
		err     error
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
		{
			number: 0,
			sides:  4,
		},
		{
			number: -3,
			sides:  6,
			err:    ErrInvalidNumberOfDice,
		},
		{
			number: 3,
			sides:  -6,
			err:    ErrInvalidNumberOfSides,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d) number %d, sides %d", i, tc.number, tc.sides), func(t *testing.T) {
			rolls, max, err := RollMax(tc.number, tc.sides)
			if len(rolls) != tc.rollLen {
				t.Errorf("[len] want %d, got %d", tc.rollLen, len(rolls))
			}

			wantMax := 0
			for r, roll := range rolls {
				if r == 0 {
					wantMax = roll
					continue
				}
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

			if err != tc.err {
				t.Errorf("[err] want %s, got %s", tc.err, err)
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
		err     error
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
		{
			number: 0,
			sides:  100,
		},
		{
			number: -1,
			sides:  6,
			err:    ErrInvalidNumberOfDice,
		},
		{
			number: 2,
			sides:  -4,
			err:    ErrInvalidNumberOfSides,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d) number %d, sides %d", i, tc.number, tc.sides), func(t *testing.T) {
			rolls, min, err := RollMin(tc.number, tc.sides)
			if len(rolls) != tc.rollLen {
				t.Errorf("[len] want %d, got %d", tc.rollLen, len(rolls))
			}

			wantMin := 0
			for r, roll := range rolls {
				if r == 0 {
					wantMin = roll
					continue
				}
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

			if err != tc.err {
				t.Errorf("[err] want %s, got %s", tc.err, err)
			}
		})
	}
}

func Test_min(t *testing.T) {
	testCases := []struct {
		want int
		num1 int
		num2 int
	}{
		{
			want: 3,
			num1: 7,
			num2: 3,
		},
		{
			want: -2,
			num1: -2,
			num2: 0,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d)", i), func(t *testing.T) {
			got := min(tc.num1, tc.num2)
			if got != tc.want {
				t.Errorf("want %d, got %d", tc.want, got)
			}
		})
	}
}

func Test_max(t *testing.T) {
	testCases := []struct {
		want int
		num1 int
		num2 int
	}{
		{
			want: 7,
			num1: 7,
			num2: 3,
		},
		{
			want: 0,
			num1: -2,
			num2: 0,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d)", i), func(t *testing.T) {
			got := max(tc.num1, tc.num2)
			if got != tc.want {
				t.Errorf("want %d, got %d", tc.want, got)
			}
		})
	}
}
