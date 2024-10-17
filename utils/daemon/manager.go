package daemon

import (
	xzfKafka "prettyy-server-online/custom-pkg/xzf-kafka"
	kafkaConsumer "prettyy-server-online/custom-pkg/xzf-kafka/consumer"
)

// KafkaJob 用于消费kafka的离线任务
// 参数 name,group: 用于读取kafka的配置文件
// 参数 consumer: 处理该kafka消息的具体逻辑
func KafkaJob(name string, group string, consumer kafkaConsumer.Actuator) func(closed <-chan struct{}) {
	return func(closed <-chan struct{}) {
		c := kafkaConsumer.NewConsumer(xzfKafka.NewConsumerConfigByNameGroup(name, group), consumer)
		go func() {
			<-closed // 容器关闭后，停止消费
			_ = c.Close()

		}()
		c.Consume()
		// 如果 consumer 实现了 Wait() 函数，则认为 consumer 需要等待 worker 退出
		if wait, ok := consumer.(interface{ Wait() }); ok {
			wait.Wait()
		}
	}
}
