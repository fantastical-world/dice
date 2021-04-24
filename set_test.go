package dice

import (
	"strings"
	"testing"
)

func TestAddingToSet(t *testing.T) {
	basicSet := Set{Name: "my dice"}
	err := basicSet.AddDice("main weapon", "1d20+3")
	if err != nil {
		t.Errorf("error encountered when adding dice %s", err)
	}
	rolls, result, err := basicSet.RollDice("main weapon")
	if err != nil {
		t.Errorf("error encountered when rolling %s", err)
	}
	//only one roll should be present
	if len(rolls) != 1 {
		t.Errorf("expected 1 roll result, but found %d results\n", len(rolls))
	}
	//single roll means sum and roll should be equal
	if (rolls[0] + 3) != result {
		t.Errorf("expected roll result and sum to be equal, but roll was %d and sum was %d\n", rolls[0], result)
	}
	//roll should be 1, 2, 3, 4, 5, ... 20
	if rolls[0] < 1 || rolls[0] > 20 {
		t.Errorf("expected roll result to be 1, 2, 3, 4, 5, ... 20, but roll was %d\n", rolls[0])
	}
}

func TestAddingDiceWithBadExpression(t *testing.T) {
	basicSet := Set{Name: "my dice"}
	err := basicSet.AddDice("main weapon", "foobar")
	if err == nil {
		t.Errorf("expected an error, but none was received")
	}
}

func TestRollingWithoutDice(t *testing.T) {
	basicSet := Set{Name: "my dice"}
	_, _, err := basicSet.RollDice("main weapon")
	if err == nil {
		t.Errorf("expected an error, but none was received")
	}
}

func TestRollingDiceThatDoesNotExist(t *testing.T) {
	basicSet := Set{Name: "my dice"}
	err := basicSet.AddDice("main weapon", "1d20+3")
	if err != nil {
		t.Errorf("error encountered when adding dice %s", err)
	}
	_, _, err = basicSet.RollDice("not here")
	if err == nil {
		t.Errorf("expected an error, but none was received")
	}
}

func TestList(t *testing.T) {
	basicSet := Set{Name: "my dice"}
	err := basicSet.AddDice("main weapon", "1d20+3")
	if err != nil {
		t.Errorf("error encountered when adding dice %s", err)
	}
	err = basicSet.AddDice("secondary weapon", "3d6")
	if err != nil {
		t.Errorf("error encountered when adding dice %s", err)
	}
	err = basicSet.AddDice("Dex Save", "1d20+4")
	if err != nil {
		t.Errorf("error encountered when adding dice %s", err)
	}

	listing := basicSet.ListDice()

	sb := strings.Builder{}
	sb.WriteString("Dex Save 1d20+4\n")
	sb.WriteString("main weapon 1d20+3\n")
	sb.WriteString("secondary weapon 3d6\n")

	expected := sb.String()

	if expected != listing {
		t.Errorf("listing not as expected\nexpected:\n%s\nactual:\n%s\n", expected, listing)
	}
}

func TestNoDiceList(t *testing.T) {
	basicSet := Set{Name: "my dice"}
	listing := basicSet.ListDice()

	if listing != "no dice" {
		t.Errorf("listing not as expected\nexpected:\nno dice\nactual:\n%s\n", listing)
	}
}
