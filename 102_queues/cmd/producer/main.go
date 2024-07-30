package producer

import (
	"strings"

	"github.com/todaatsushi/queue/internal/producer"
)

func Send(port int, messages string) {
	split := strings.Split(messages, ",")
	producer.QueueTasks(port, split...)
}

func Health(port int) {
	producer.CheckServer(port)
}

func QueueLen(port int) {
	producer.GetQueueLen(port)
}
