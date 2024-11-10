import abc
import dataclasses as dc


@dc.dataclass
class BaseServerConfig(abc.ABC):
    host: str
    port: int


@dc.dataclass
class BasicServer(BaseServerConfig):
    pass


# Replace with env / file config
SERVERS = (
    BasicServer(host="localhost", port=8000),
    BasicServer(host="localhost", port=8001),
    BasicServer(host="localhost", port=8002),
    BasicServer(host="localhost", port=8003),
)
