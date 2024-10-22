package metrics

var CommonCounter *commonCounter

type commonCounter struct {
	counter *counter
}

func (acc *commonCounter) Inc(group, name string) {
	acc.counter.Values(group, name).Inc()
}

func newCommonCounter() *commonCounter {
	acc := &commonCounter{
		counter: newCounter("common", "common counter", []string{"group", "name"}),
	}
	return acc
}

func init() {
	CommonCounter = newCommonCounter()
}
