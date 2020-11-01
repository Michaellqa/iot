package aggregation_test

import (
	"github.com/Michaellqa/iot/internal/aggregation"
	"github.com/Michaellqa/iot/internal/aggregation/mock_aggregation"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

var writeDuration = 100 * time.Millisecond

func TestAsyncStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mock_aggregation.NewMockStorage(ctrl)
	fifo := mock_aggregation.NewMockFifo(ctrl)
	as := aggregation.NewAsyncStorage(fifo, store)

	cases := []struct {
		name               string
		setupExpectedCalls func()
	}{
		{
			name: "",
			setupExpectedCalls: func() {
				record := aggregation.Record{Id: "id_1", Value: float64(1)}
				fifo.EXPECT().Add(record)
				fifo.EXPECT().Get().Return(interface{}(record))
				store.EXPECT().Write(record).Do(func(interface{}) {
					time.Sleep(writeDuration)
				})

				as.Add(record)

				fifo.EXPECT().Len()

				as.Wait()
			},
		},
	}
	for _, c := range cases {
		c.setupExpectedCalls()
	}
}
