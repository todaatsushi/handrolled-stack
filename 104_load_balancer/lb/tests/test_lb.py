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


def create_strategy(
    klass: type[strategies.Strategy], num_servers: int
) -> strategies.Strategy:
    _servers = tuple(
        servers.BasicServer(host="local", port=1 + i) for i in range(num_servers)
    )
    return klass.new(_servers)


class TestForwarding(unittest.TestCase):
    def setUp(self) -> None:
        app = create_app()
        self.client = get_client(app)

    @mock.patch("lb.servers.requests.request")
    @mock.patch("lb.strategies.get_strategy")
    def test_forwards_to_server(
        self, mock_get_strategy: mock.Mock, mock_request: mock.Mock
    ) -> None:
        mock_get_strategy.return_value = create_strategy(strategies.RoundRobin, 1)
        mock_request.return_value = mock.Mock(
            status_code=200, content=b"OK", headers={}, spec=requests.Response
        )

        response = self.client.get("/")
        self.assertEqual(response.status_code, 200)

    @mock.patch("lb.servers.requests.request")
    @mock.patch("lb.strategies.get_strategy")
    def test_forwards_to_next_if_bad_response(
        self, mock_get_strategy: mock.Mock, mock_request: mock.Mock
    ) -> None:
        mock_get_strategy.return_value = create_strategy(strategies.RoundRobin, 1)
        mock_request.side_effect = [
            mock.Mock(status_code=503, content=b"", headers={}, spec=requests.Response),
            mock.Mock(
                status_code=200, content=b"ok", headers={}, spec=requests.Response
            ),
        ]

        response = self.client.get("/")
        self.assertEqual(response.status_code, 200)
