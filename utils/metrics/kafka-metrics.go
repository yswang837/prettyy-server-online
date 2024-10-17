package metrics

var KafkaCounter *kafkaCounter

const (
	Success     = "success"
	Timeout     = "timeout"
	ErrSend     = "err_send"
	ErrDatatype = "err_datatype"
)

type KafkaCounterTags struct {
	Name string
	Type string
}

type kafkaCounter struct {
	counter *counter
}

func (acc *kafkaCounter) Inc(tags *KafkaCounterTags) {
	acc.counter.Values(tags.Name, tags.Type).Inc()
}

func newKafkaCounter() *kafkaCounter {
	acc := &kafkaCounter{
		counter: newCounter("kafka", "kafka counter", []string{"name", "type"}),
	}
	return acc
}

func init() {
	KafkaCounter = newKafkaCounter()
}
