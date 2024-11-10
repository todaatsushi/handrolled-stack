import flask

import uuid
from lb import strategies


# Env var
MAX_ATTEMPS = 2
strategy = strategies.get_strategy()


def configure_routes(app: flask.Flask) -> None:
    @app.route("/")
    def lb() -> flask.Response:
        request_id = uuid.uuid4()

        tries = 0
        while tries < MAX_ATTEMPS:
            dest = strategy.get_next()
            response = dest.forward(flask.request, str(request_id))
            if response.status_code >= 500:
                tries += 1
                continue
            return response
        return response
