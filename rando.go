package dice

import "math/rand"

type seeder struct {
	seed   int64
	random *rand.Rand
}

func (s *seeder) Seed(seed int64) {
	s.seed = seed
}

func (s *seeder) Reset() {
	s.random = rand.New(rand.NewSource(s.seed))
}

func (s *seeder) RandomRange(min, max int) int {
	if min > max {
		return 0
	}

	if min == max {
		return min
	}

	return s.random.Intn((max+1)-min) + min
}

func (s *seeder) RandomNRange(count, min, max int, unique bool) []int {
	if count < 1 {
		return nil
	}

	var results []int

	if min > max {
		if unique {
			return []int{0}
		}
		for i := 0; i < count; i++ {
			results = append(results, 0)
		}
		return results
	}

	if min == max {
		if unique {
			return []int{min}
		}
		for i := 0; i < count; i++ {
			results = append(results, min)
		}
		return results
	}

	if unique {
		// this adjusts count if there are not enough unique values in the provided range
		if count > ((max + 1) - min) {
			count = (max + 1) - min
		}

		temp := make(map[int]int)

		for len(temp) < count {
			v := s.RandomRange(min, max)
			_, exists := temp[v]
			if !exists {
				temp[v] = v
				results = append(results, v)
			}
		}

		return results
	}

	for i := 0; i < count; i++ {
		results = append(results, s.RandomRange(min, max))
	}

	return results
}

func New(seed int64) *seeder {
	s := &seeder{seed: seed}
	s.Reset()

	return s
}
