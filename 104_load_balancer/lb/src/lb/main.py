import flask


app = flask.Flask(__name__)


@app.route("/")
def lb() -> str:
    return "LB\n"
