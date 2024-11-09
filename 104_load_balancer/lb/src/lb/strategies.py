from __future__ import annotations

import collections
from typing import Final, Protocol, TypeVar

from lb import servers

T = TypeVar("T")


class Strategy(Protocol[T]):
    @classmethod
    def new(cls, servers: tuple[T, ...]) -> Strategy[T]:
        raise NotImplementedError

    def get_next(self) -> T:
        raise NotImplementedError


class RoundRobin(collections.deque, Strategy[servers.BaseServerConfig]):
    tag: Final = "ROUND_ROBIN"

    @classmethod
    def new(
        cls,
        servers: tuple[servers.BaseServerConfig, ...],
    ) -> "RoundRobin":
        queue = cls()
        for server in servers:
            queue.append(server)
        return queue

    def get_next(self) -> str:
        _next = self.popleft()
        self.append(_next)
        return _next


STRATEGY_MAP = {
    RoundRobin.tag: RoundRobin,
}
