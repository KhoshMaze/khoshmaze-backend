from config.config import Config
from app.redis.redis import Redis

class App:

    __db = None
    def __init__(self, config: Config, debug: bool = False):
        self.__config = config
        self.__db == self.__set_db()
        self.DEBUG = debug # CONSTANT VALUE. DO NOT USE IT ANYWHERE

    def __set_db(self): 
        r = Redis(self.__config.redis)
        if r is None:
            raise Exception("Redis server is not running")
        else:
            return r

        
    def DB(self):
        if self.__db is None:
            self.__db = self.__set_db()
        return self.__db
    
    def ServerConfig(self):
        return self.__config.server

    def BotConfig(self):
        return self.__config.bot
