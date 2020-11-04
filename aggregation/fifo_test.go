package aggregation

import "testing"

func TestList_Add_Get(t *testing.T) {

	cases := []struct {
		name        string
		sequence    func(q *ListFifo)
		expectedSum int
	}{
		{
			name:        "no additions",
			sequence:    func(fifo *ListFifo) {},
			expectedSum: 0,
		},
		{
			name: "additions only",
			sequence: func(q *ListFifo) {
				q.Add(1)
				q.Add(2)
				q.Add(4)
			},
			expectedSum: 7,
		},
		{
			name: "mixed adds and gets",
			sequence: func(q *ListFifo) {
				q.Get()
				q.Add(1)
				q.Add(2)
				q.Get()
				q.Add(4)
				q.Get()
				q.Add(8)
			},
			expectedSum: 12,
		},
	}
	for _, c := range cases {
		q := &ListFifo{}
		c.sequence(q)
		res := queueSum(q)

		if res != c.expectedSum {
			t.Errorf("%s\nGot  %v\nWant %v\n", c.name, res, c.expectedSum)
		}
	}
}

func queueSum(q *ListFifo) int {
	res := 0
	for {
		v := q.Get()
		if v == nil {
			break
		}
		res += v.(int)
	}
	return res
}
