import os 
import sys
import importlib.util

_package_dir = os.path.dirname(os.path.abspath(__file__))

# Add this directory to sys.path temporarily
if _package_dir not in sys.path:
    sys.path.insert(0, _package_dir)

import food_pb2
import food_pb2_grpc
import common_pb2
import common_pb2_grpc

food_pb2 = food_pb2
food_pb2_grpc = food_pb2_grpc
    
if _package_dir in sys.path:
    sys.path.remove(_package_dir)