import unittest

from lb import servers, strategies


def create_servers() -> tuple[servers.BaseServerConfig, ...]:
    return tuple(servers.BasicServer(host="local", port=1 + i) for i in range(3))


class TestRoundRobin(unittest.TestCase):
    def test_get_next(self) -> None:
        all_servers = create_servers()
        rr = strategies.RoundRobin.new(all_servers)

        expected = [
            all_servers[0],
            all_servers[1],
            all_servers[2],
        ]
        actual = [rr.get_next() for _ in range(3)]

        self.assertEqual(expected, actual)

    def test_readded_to_pool(self) -> None:
        all_servers = create_servers()
        rr = strategies.RoundRobin.new(all_servers)

        self.assertEqual(len(rr), 3)

        first = rr.get_next()
        self.assertEqual(len(rr), 3)
        last = rr.pop()

        self.assertEqual(first, last)
