package consumer

import (
	cluster "github.com/bsm/sarama-cluster"
	"log"
	xzfKafka "prettyy-server-online/custom-pkg/xzf-kafka"
	"sync"
)

type simpleConsumerGroup struct {
	*cluster.Consumer
	actuator Actuator
	wg       sync.WaitGroup
	debug    bool
}

func (s *simpleConsumerGroup) init(conf *xzfKafka.Config, act Actuator) error {
	clusterConf := &cluster.Config{Config: *conf.Config}
	clusterConf.Group.PartitionStrategy = cluster.Strategy(conf.Config.Consumer.Group.Rebalance.Strategy.Name())
	clusterConf.Group.Mode = cluster.ConsumerMode(cluster.ConsumerModePartitions)
	clusterConf.Group.Offsets.Retry.Max = conf.Config.Consumer.Offsets.Retry.Max
	clusterConf.Group.Offsets.Synchronization.DwellTime = conf.Consumer.MaxProcessingTime
	clusterConf.Group.Session.Timeout = conf.Consumer.Group.Session.Timeout
	clusterConf.Group.Heartbeat.Interval = conf.Consumer.Group.Heartbeat.Interval
	clusterConf.Config.Version = conf.Version
	clusterConf.Consumer.Offsets.Initial = conf.Consumer.Offsets.Initial

	consumer, err := cluster.NewConsumer(conf.Brokers, conf.GroupId, conf.Topics, clusterConf)
	if err != nil {
		return err
	}

	if err := act.Init(); err != nil {
		return err
	}
	s.Consumer = consumer
	s.actuator = act
	s.debug = conf.Debug
	return nil
}

func (s *simpleConsumerGroup) consume() {
	for part := range s.Consumer.Partitions() {
		s.wg.Add(1)
		go func(pc cluster.PartitionConsumer) {
			for message := range pc.Messages() {
				if err := s.actuator.Do(message); err != nil {
					log.Printf("Message claimed err: topic = %s, partition = %d, offset = %d, key = %s, value = %s, err = %s", message.Topic, message.Partition, message.Offset, message.Key, string(message.Value), err)

				}
				if !s.debug {
					s.Consumer.MarkOffset(message, "")
				}
			}
			s.wg.Done()
		}(part)
	}
}

func (s *simpleConsumerGroup) close() error {
	return s.Consumer.Close()
}
