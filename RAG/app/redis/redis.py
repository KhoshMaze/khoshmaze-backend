import redis
 
from config.config import RedisConfig

class Redis(redis.Redis):

    KEY_PREFIX = "RAG"
    def __init__(self, config: RedisConfig):
        super().__init__(host=config.host, port=config.port, db=config.db,password=config.password)
        
    def get(self, key : str):
        return super().get(self.KEY_PREFIX + '.' + key)
    
    def set(self, key : str, value : str, ttl: int = 120):
        return super().set(self.KEY_PREFIX + '.' + key, value, ex=ttl)
    
    def delete(self, key : str):
        return super().delete(self.KEY_PREFIX + '.' + key)

    def exists(self, key : str):
        return super().exists(self.KEY_PREFIX + '.' + key)
    
    def update_ttl(self, key : str, ttl: int):
        if self.exists(key):
            super().expire(self.KEY_PREFIX + '.' + key, ttl)
    
    def get_all_keys(self):
        return super().keys(self.KEY_PREFIX + '.*')

