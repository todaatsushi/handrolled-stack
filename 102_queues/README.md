# Task queue

Create a message broker software which can accept and produce messages sent to it.

## Requirements
- [ ] Send messages, add to queue
- [ ] Send consume message, remove from queue and return
- [ ] Multiple consumers should be able to read off the queue
- [ ] Check number of items on the queue

# Extra
- [ ] FIFO - lock the consumers until tasks have finished running
- [ ] Sending strategies e.g. retry with backoff(at least once etc.).
- [ ] Pub sub - notify when tasks added to queue
