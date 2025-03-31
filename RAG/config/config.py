import json


class ServerConfig:
    def __init__(self, config: dict):
        self.port = config.get("port", 8000)
        self.host = config.get("host", "0.0.0.0")

        try:
            self.auth_secret = config.get("authSecret")
        except KeyError:
            raise ValueError("authSecret is required")


class RedisConfig:
    def __init__(self, config: dict):
        self.host = config.get("host", "localhost")
        self.port = config.get("port", 6379)
        self.password = config.get("password", "")


class Config:

    instance = None

    """    Singleton Class     """

    def __new__(cls, *args, **kwargs):
        if cls.instance is None:
            cls.instance = super().__new__(cls)
        return cls.instance

    def __init__(self, config_path: str):
        s, r = self.__load_config(config_path)
        self.__server = s
        self.__redis = r
        self.instance = self

    @staticmethod
    def __load_config(config_path: str) -> tuple[ServerConfig, RedisConfig]:
        with open(config_path, "r") as f:
            config = json.load(f)
            return ServerConfig(config["server"]), RedisConfig(config["redis"])

    def Server(self):
        return self.__server

    def Redis(self):
        return self.__redis
