# Message queues

[Message queues](https://sudhir.io/the-big-little-guide-to-message-queues)

## Spec
Create a message broker which can accept messages, queue them, and deliver them to subscribers.

It should be able to queue them up.

Also create a worker that can process the messages - doesn't have to do anything special, just 
take the messages off the queue.

Queue should be able to do:
- FIFO
- Any ordering
