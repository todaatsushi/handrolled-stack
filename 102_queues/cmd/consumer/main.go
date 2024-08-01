package consumer

import "github.com/todaatsushi/queue/internal/consumer"

func Start(port int, numConsumers int) {
	consumer.StartConsumers(port, numConsumers)
}
