from __future__ import annotations
from typing import TYPE_CHECKING

import unittest
from unittest import mock
import flask
import requests

from lb import servers, strategies
from lb.app import app as main

if TYPE_CHECKING:
    from flask.testing import FlaskClient


def create_app() -> flask.Flask:
    app = main.create_app()
    app.config.update({"TESTING": True})
    return app


def get_client(app: flask.Flask | None = None) -> FlaskClient:
    if app is None:
        app = create_app()
    return app.test_client()


class TestForwarding(unittest.TestCase):
    def setUp(self) -> None:
        app = create_app()
        self.client = get_client(app)

        _servers = (servers.BasicServer(host="test", port=8000),)
        self.strategy = strategies.RoundRobin.new(_servers)

    @mock.patch("lb.servers.requests.request")
    @mock.patch("lb.strategies.get_strategy")
    def test_forwards_to_server(
        self, mock_get_strategy: mock.Mock, mock_request: mock.Mock
    ) -> None:
        mock_get_strategy.return_value = self.strategy
        mock_request.return_value = mock.Mock(
            status_code=200, content=b"OK", headers={}, spec=requests.Response
        )

        response = self.client.get("/")
        self.assertEqual(response.status_code, 200)
