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
SERVERS = (BasicServer(host="localhost", port=5000),)
