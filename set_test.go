package dice

import (
	"reflect"
	"strings"
	"testing"
)

func TestSet_AddDice(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		want := Set{
			Name: "my dice",
			Dice: map[string]string{
				"main weapon": "1d20+3",
			},
		}

		got := Set{Name: "my dice"}
		err := got.AddDice("main weapon", "1d20+3")
		if err != nil {
			t.Errorf("unexpected error, %s", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("error when expression is invalid", func(t *testing.T) {
		want := Set{
			Name: "my dice",
		}

		got := Set{Name: "my dice"}
		err := got.AddDice("main weapon", "hey0d20+2")
		if err != ErrInvalidRollExpression {
			t.Errorf("want %s, got %s", ErrInvalidRollExpression, err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
	})
}

func TestSet_RollDice(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		subject := Set{Name: "my dice"}
		err := subject.AddDice("main weapon", "1d20+3")
		if err != nil {
			t.Errorf("unexpected error, %s", err)
		}

		rolls, sum, err := subject.RollDice("main weapon")
		if err != nil {
			t.Errorf("unexpected error, %s", err)
		}
		//simple test only since roll tests cover much more
		if len(rolls) != 1 {
			t.Errorf("[len] want %d, got %d", 1, len(rolls))
		}
		wantSum := 0
		for _, roll := range rolls {
			wantSum += roll
		}
		wantSum += 3
		if wantSum != sum {
			t.Errorf("[sum] want %d, got %d", wantSum, sum)
		}

		for _, roll := range rolls {
			if roll < 1 || roll > 20 {
				t.Errorf("[rolls] want roll %d-%d, got %d", 1, 20, roll)
			}
		}
	})

	t.Run("error when no dice in set", func(t *testing.T) {
		subject := Set{Name: "my dice"}
		_, _, err := subject.RollDice("main weapon")
		if err != ErrEmptyDiceSet {
			t.Errorf("want %s, got %s", ErrEmptyDiceSet, err)
		}
	})

	t.Run("error when specified dice does not exist", func(t *testing.T) {
		subject := Set{Name: "my dice"}
		err := subject.AddDice("main weapon", "1d20+3")
		if err != nil {
			t.Errorf("unexpected error, %s", err)
		}

		_, _, err = subject.RollDice("no dice")
		if err != ErrDiceNotFound {
			t.Errorf("want %s, got %s", ErrDiceNotFound, err)
		}
	})
}

func TestSet_ListDice(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		sb := strings.Builder{}
		sb.WriteString("Dex Save 1d20+4\n")
		sb.WriteString("main weapon 1d20+3\n")
		sb.WriteString("secondary weapon 3d6\n")
		want := sb.String()

		subject := Set{Name: "my dice"}
		err := subject.AddDice("main weapon", "1d20+3")
		if err != nil {
			t.Errorf("unexpected error, %s", err)
		}
		err = subject.AddDice("secondary weapon", "3d6")
		if err != nil {
			t.Errorf("unexpected error, %s", err)
		}
		err = subject.AddDice("Dex Save", "1d20+4")
		if err != nil {
			t.Errorf("unexpected error, %s", err)
		}

		got := subject.ListDice()
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("empty set", func(t *testing.T) {
		want := ""
		subject := Set{Name: "my dice"}
		got := subject.ListDice()
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})
}
