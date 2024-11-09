import abc
import dataclasses as dc


@dc.dataclass
class BaseServerConfig(abc.ABC):
    host: str
    port: int


@dc.dataclass
class BasicServer(BaseServerConfig):
    pass
