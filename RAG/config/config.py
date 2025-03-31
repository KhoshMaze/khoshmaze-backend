import json


class ServerConfig:
    def __init__(self, config: dict):
        self.port = config.get("port", 8000)
        self.host = config.get("host", "0.0.0.0")

        try:
            self.auth_secret = config["authSecret"]
        except KeyError:
            raise ValueError("authSecret is required")

        try:
            self.grpc_api_key = config["grpc_api_key"]
        except KeyError:
            raise ValueError("grpc_api_key is required")


class RedisConfig:
    def __init__(self, config: dict):
        self.host = config.get("host", "localhost")
        self.port = config.get("port", 6379)
        self.password = config.get("password", "")
        self.db = config.get("db", 0)


class AIConfig:
    def __init__(self, config: dict):
        self.open_api_key = config.get("open_api_key", "")
        self.tracing = config.get("tracing", False)


class Config:

    instance = None

    """    Singleton Class     """

    def __new__(cls, *args, **kwargs):
        if cls.instance is None:
            cls.instance = super().__new__(cls)
        return cls.instance

    def __init__(self, config_path: str):
        s, r, b = self.__load_config(config_path)
        self.__server = s
        self.__redis = r
        self.__bot = b
        self.instance = self
        
    @staticmethod
    def __load_config(config_path: str) -> tuple[ServerConfig, RedisConfig, AIConfig]:
        with open(config_path, "r") as f:
            config = json.load(f)
            return (
                ServerConfig(config["server"]),
                RedisConfig(config["redis"]),
                AIConfig(config["bot"]),
            )

    def Server(self):
        return self.__server

    def Redis(self):
        return self.__redis

    def Bot(self):
        return self.__bot