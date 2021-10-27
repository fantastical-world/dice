package dice

import (
	"fmt"
	"testing"
)

func Test_RollChallenge(t *testing.T) {
	testCases := []struct {
		expression    string
		against       int
		equalSucceeds bool
		alert         []int
		err           error
	}{
		{
			expression:    "1d20+3",
			against:       13,
			equalSucceeds: false,
			alert:         nil,
			err:           nil,
		},
		{
			expression:    "1d20+3",
			against:       13,
			equalSucceeds: true,
			alert:         nil,
			err:           nil,
		},
		{
			expression:    "1d1+3",
			against:       4,
			equalSucceeds: false,
			alert:         nil,
			err:           nil,
		},
		{
			expression:    "1d1+3",
			against:       4,
			equalSucceeds: true,
			alert:         nil,
			err:           nil,
		},
		{
			expression:    "1d2+3",
			against:       5,
			equalSucceeds: true,
			alert:         []int{1, 2},
			err:           nil,
		},
		{
			expression:    "hey0+3",
			against:       5,
			equalSucceeds: true,
			alert:         []int{1, 2},
			err:           ErrInvalidRollExpression,
		},
		{
			expression:    "-2d6",
			against:       5,
			equalSucceeds: true,
			alert:         nil,
			err:           ErrInvalidRollExpression,
		},
		{
			expression:    "2d-20",
			against:       5,
			equalSucceeds: true,
			alert:         nil,
			err:           ErrInvalidRollExpression,
		},
		{
			expression:    "3d4+2",
			against:       8,
			equalSucceeds: false,
			alert:         []int{1, 2, 3, 4},
			err:           nil,
		},
		{
			expression:    "0d4+2",
			against:       2,
			equalSucceeds: false,
			alert:         []int{2},
			err:           nil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d) %s", i, tc.expression), func(t *testing.T) {
			succeeded, result, found, err := RollChallenge(tc.expression, tc.against, tc.equalSucceeds, tc.alert)

			if tc.equalSucceeds {
				if succeeded && (result < tc.against) {
					t.Errorf("succeded even though %d is not >= %d", result, tc.against)
				}
			} else {
				if succeeded && (result <= tc.against) {
					t.Errorf("succeded even though %d is not > %d", result, tc.against)
				}
			}

			if tc.alert == nil && (len(found) > 0) {
				t.Errorf("[alert] wanted no alert, got %v", found)
			}

			for _, f := range found {
				requestedAlert := false
				for _, a := range tc.alert {
					if f == a {
						requestedAlert = true
					}
				}

				if !requestedAlert {
					t.Errorf("[alert] wanted alert on %v, got alert on %d", tc.alert, f)
				}
			}

			if err != tc.err {
				t.Errorf("[err] want %s, got %s", tc.err, err)
			}
		})
	}
}
