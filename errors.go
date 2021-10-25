package dice

type DiceError string

func (d DiceError) Error() string { return string(d) }

var (
	ErrInvalidRollExpression = DiceError("not a valid roll expression")
	ErrEmptyDiceSet          = DiceError("you do not have any dice in your set")
	ErrDiceNotFound          = DiceError("dice not found")
	ErrInvalidOperator       = DiceError("invalid operator")
)
