import logging
import flask


logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)


def create_app() -> flask.Flask:
    app = flask.Flask(__name__)

    from .routes import configure_routes

    with app.app_context():
        configure_routes(app)
    return app


app = create_app()
