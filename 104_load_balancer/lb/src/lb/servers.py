import abc
import dataclasses as dc
import uuid
import requests

import flask


@dc.dataclass
class BaseServerConfig(abc.ABC):
    host: str
    port: int

    def forward(
        self, request: flask.Request, request_id: str | None = None
    ) -> flask.Response:
        forward_addr = f"http://{self.host}:{self.port}{request.path}"

        request_headers = dict(request.headers)
        request_id = request_id or str(uuid.uuid4())
        request_headers["X-Request-ID"] = str(request_id)

        forward = requests.request(
            method=request.method,
            url=forward_addr,
            headers=request_headers,
            data=request.data,
            cookies=request.cookies,
            allow_redirects=False,
        )
        return flask.Response(
            forward.content, forward.status_code, forward.headers.items()
        )


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
