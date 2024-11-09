from __future__ import annotations

import logging
import flask

logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)

app = flask.Flask(__name__)
SERVER_ID = 1


@app.route("/")
def home() -> str:
    logger.info(
        f"Server #{SERVER_ID} received a request: {flask.request.method} @ {flask.request.url}"
    )
    return f"Hello from server #{SERVER_ID}\n"
