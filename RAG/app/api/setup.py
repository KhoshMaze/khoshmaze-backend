import uvicorn
from fastapi import FastAPI
from app.app import App


class Server:
    def __init__(self, appContainer: App):
        self.__appContainer = appContainer
        self.__server = FastAPI()
        self.__register_routes()

    def __register_routes(self):
        server = self.__server

    def run(self):
        config = self.__appContainer.ServerConfig()
        uvicorn.run(self.__server, host=config.host, port=config.port)
