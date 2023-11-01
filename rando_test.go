package dice

import (
	"math/rand"
	"reflect"
	"testing"
)

func Test_seeder_Seed(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		want := int64(456)
		subject := &seeder{seed: 123}
		subject.Seed(want)
		if subject.seed != want {
			t.Errorf("want %d, got %d", want, subject.seed)
		}
	})
}

func Test_seeder_Reset(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		subject := &seeder{}
		if subject.random != nil {
			t.Error("expected nil")
		}
		subject.Reset()
		if subject.random == nil {
			t.Error("expected non-nil")
		}
	})
}

func Test_seeder_RandomRange(t *testing.T) {
	testCases := []struct {
		scenario string
		seed     int64
		min      int
		max      int
		want     int
	}{
		{
			scenario: "min greater than max",
			seed:     100,
			min:      30,
			max:      10,
			want:     0,
		},
		{
			scenario: "min and max are equal",
			seed:     100,
			min:      99,
			max:      99,
			want:     99,
		},
		{
			scenario: "a value from 1 to 20",
			seed:     100,
			min:      1,
			max:      20,
			want:     4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			subject := &seeder{random: rand.New(rand.NewSource(tc.seed))}
			got := subject.RandomRange(tc.min, tc.max)
			if got != tc.want {
				t.Errorf("want %d, got %d", tc.want, got)
			}
		})
	}
}

func Test_seeder_RandomNRange(t *testing.T) {
	testCases := []struct {
		scenario string
		seed     int64
		count    int
		min      int
		max      int
		unique   bool
		want     []int
	}{
		{
			scenario: "want 0",
			seed:     100,
			count:    0,
			min:      1,
			max:      10,
			unique:   false,
		},
		{
			scenario: "min greater than max - want 2",
			seed:     100,
			count:    2,
			min:      30,
			max:      10,
			unique:   false,
			want:     []int{0, 0},
		},
		{
			scenario: "min and max are equal - want 2",
			seed:     100,
			count:    2,
			min:      99,
			max:      99,
			unique:   false,
			want:     []int{99, 99},
		},
		{
			scenario: "min greater than max - want 2 unique",
			seed:     100,
			count:    2,
			min:      30,
			max:      10,
			unique:   true,
			want:     []int{0},
		},
		{
			scenario: "min and max are equal - want 2 unique",
			seed:     100,
			count:    2,
			min:      99,
			max:      99,
			unique:   true,
			want:     []int{99},
		},
		{
			scenario: "a value from 1 to 20",
			seed:     100,
			count:    1,
			min:      1,
			max:      20,
			unique:   false,
			want:     []int{4},
		},
		{
			scenario: "a value from 1 to 20 - want 5",
			seed:     99,
			count:    5,
			min:      1,
			max:      20,
			unique:   false,
			want:     []int{18, 4, 11, 3, 3},
		},
		{
			scenario: "a value from 1 to 20 - want 5 unique",
			seed:     99,
			count:    5,
			min:      1,
			max:      20,
			unique:   true,
			want:     []int{18, 4, 11, 3, 2},
		},
		{
			scenario: "a value from 1 to 5 - want 10",
			seed:     99,
			count:    10,
			min:      1,
			max:      5,
			unique:   false,
			want:     []int{3, 4, 1, 3, 3, 4, 2, 5, 5, 3},
		},
		{
			scenario: "a value from 1 to 5 - want 10 unique (will only have 5)",
			seed:     99,
			count:    10,
			min:      1,
			max:      5,
			unique:   true,
			want:     []int{3, 4, 1, 2, 5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			subject := &seeder{random: rand.New(rand.NewSource(tc.seed))}
			got := subject.RandomNRange(tc.count, tc.min, tc.max, tc.unique)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}

func Test_New(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		subject := New(100)
		if subject == nil {
			t.Error("expected non-nil")
		}
	})
}
