import sys
from config.config import Config
from app.app import App
from app.api.setup import Server

def main():

    if "--config" in sys.argv:
        path = sys.argv[sys.argv.index("--config") + 1]
    else:
        path = "./config.json"

    config = Config(path)

    if "--debug" in sys.argv:
        debug = True
    else:
        debug = False

    app = App(config, debug)

    server = Server(app)
    server.run()
    
if __name__ == "__main__":
    main()
