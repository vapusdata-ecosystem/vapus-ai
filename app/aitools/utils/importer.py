import os
import sys

def proto_importer():
    BASE_DIR = os.path.dirname(os.path.abspath(__file__))
    PROTOBUF_DIR = os.path.join(BASE_DIR, '../../../apis', 'gen-python')
    PROTOBUF_DIR = os.path.abspath(PROTOBUF_DIR)
    print("DEBUG: Adding to sys.path:", PROTOBUF_DIR)
    sys.path.append(PROTOBUF_DIR)