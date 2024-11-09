from __future__ import annotations

from typing import Protocol, TypeVar

T = TypeVar("T")


class Strategy(Protocol[T]):
    @classmethod
    def new(cls, servers: tuple[T, ...]) -> Strategy[T]:
        raise NotImplementedError

    def get_next(self) -> T:
        raise NotImplementedError
