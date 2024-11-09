import logging

logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)

def run() -> None:
    logger.info("Running main function")

if __name__ == '__main__':
    run()
