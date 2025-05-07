from typing import Any

# use class and setter functions to store global variables
def set_secret_store(val: Any):
    global SECRET_VAR
    SECRET_VAR = val

def set_service_config(val: Any):
    global SERVICE_CONFIG_VAR
    SERVICE_CONFIG_VAR = val