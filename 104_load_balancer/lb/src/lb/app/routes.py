import flask

import uuid
from lb import strategies


strategy = strategies.get_strategy()


def configure_routes(app: flask.Flask) -> None:
    @app.route("/")
    def lb() -> flask.Response:
        request_id = uuid.uuid4()
        dest = strategy.get_next()
        return dest.forward(flask.request, str(request_id))
