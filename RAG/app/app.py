from config.config import Config
from app.redis.redis import Redis

class App:

    __db = None
    def __init__(self, config: Config):
        self.config = config
        self.__db == self.__set_db()

    def __set_db(self): 
        r = Redis(self.config.Redis())
        if r is None:
            raise Exception("Redis server is not running")
        else:
            return r

        
    def DB(self):
        if self.__db is None:
            self.__db = self.__set_db()
        return self.__db
