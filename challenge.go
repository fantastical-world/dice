package dice

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

	succeeded := result > against

	if !succeeded && equalSucceeds {
		succeeded = result == against
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
