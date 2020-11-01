package aggregation

import "testing"

func TestList_Add_Get(t *testing.T) {

	cases := []struct {
		name        string
		sequence    func(q *fifo)
		expectedSum int
	}{
		//{
		//	name:        "no additions",
		//	sequence:    func(q *fifo) {},
		//	expectedSum: 0,
		//},
		{
			name: "additions only",
			sequence: func(q *fifo) {
				q.add(1)
				q.add(2)
				q.add(4)
			},
			expectedSum: 7,
		},
		{
			name: "mixed adds and gets",
			sequence: func(q *fifo) {
				q.get()
				q.add(1)
				q.add(2)
				q.get()
				q.add(4)
				q.get()
				q.add(8)
			},
			expectedSum: 12,
		},
	}
	for _, c := range cases {
		q := &fifo{}
		c.sequence(q)
		res := queueSum(q)

		if res != c.expectedSum {
			t.Errorf("%s\nGot  %v\nWant %v\n", c.name, res, c.expectedSum)
		}
	}
}

func queueSum(q *fifo) int {
	res := 0
	for {
		v := q.get()
		if v == nil {
			break
		}
		res += v.(int)
	}
	return res
}
