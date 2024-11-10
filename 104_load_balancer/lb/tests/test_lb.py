import unittest
import flask


def create_app() -> flask.Flask:
    app = flask.Flask(__name__)
    app.config.update({"TESTING": True})
    return app


def get_client(app: flask.Flask | None = None) -> flask.testing.FlaskClient:
    if app is None:
        app = create_app()
    return app.test_client()


class TestForwarding(unittest.TestCase):
    def setUp(self) -> None:
        self.client = get_client()
