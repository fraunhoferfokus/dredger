package generator

type ChannelBinding struct {
	HTTP  HTTPBinding
	NATS  NATSBinding
	Kafka KafkaBinding
}
