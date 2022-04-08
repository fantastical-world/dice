package dice

type Error string

func (d Error) Error() string { return string(d) }

var (
	ErrInvalidRollExpression = Error("not a valid roll expression")
	ErrEmptyDiceSet          = Error("you do not have any dice in your set")
	ErrDiceNotFound          = Error("dice not found")
	ErrInvalidOperator       = Error("invalid operator")
	ErrInvalidNumberOfDice   = Error("invalid number of dice")
	ErrInvalidNumberOfSides  = Error("invalid number of sides")
)
