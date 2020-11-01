package generation

import "math/rand"

func NewRandomSource(
	id string,
	curValue int,
	step int,
) Source {
	return &dataSource{id: id, curValue: curValue, step: step}
}

type dataSource struct {
	id       string
	curValue int
	step     int
}

func (s *dataSource) ReadValue() Metrics {
	if s.step != 0 {
		dif := rand.Intn(2*s.step+1) - s.step
		s.curValue += dif
	}
	return Metrics{Id: s.id, Value: s.curValue}
}
