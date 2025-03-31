import sys
from config.config import Config


def main():

    if "--config" in sys.argv:
        path = sys.argv[sys.argv.index("--config") + 1]
    else:
        path = "./config.json"

    config = Config(path)


if __name__ == "__main__":
    main()
