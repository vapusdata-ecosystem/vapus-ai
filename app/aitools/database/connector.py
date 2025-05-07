import json
from sqlalchemy import create_engine
from urllib.parse import quote_plus

creds = {"secret_value":"{\"dataSourceCreds\":{\"genericCredentialObj\":{\"username\":\"postgres\",\"password\":\"\",\"aws_creds\":{},\"gcp_creds\":{},\"azure_creds\":{}},\"url\":\"35.244.50.2\",\"db\":\"vapusdata\",\"port\":5432},\"dataSourceEngine\":\"POSTGRES\",\"dataSourceType\":\"DSDT_DATABASE\",\"dataSourceService\":\"GCP_CLOUD_SQL\",\"params\":{},\"dataSourceSvcProvider\":\"SSP_GCP\"}"}

class DatabaseConnector:

    def parse_data_source_info(self,secret_value):
        
        try:
            secret_dict = json.loads(secret_value)
        except json.JSONDecodeError as e:
            raise ValueError("Invalid JSON format") from e

        
        try:
            creds = secret_dict.get("dataSourceCreds", {}).get("genericCredentialObj", {})
            url = secret_dict.get("dataSourceCreds", {}).get("url", "")
            db = secret_dict.get("dataSourceCreds", {}).get("db", "")
            port = secret_dict.get("dataSourceCreds", {}).get("port", "")
            engine = secret_dict.get("dataSourceEngine", "")
            service = secret_dict.get("dataSourceService", "")
            provider = secret_dict.get("dataSourceSvcProvider", "")

            result = {
                "username": creds.get("username", ""),
                "password": creds.get("password", ""),
                "url": url,
                "db": db,
                "port": port,
                "engine": engine,
                "service": service,
                "provider": provider,
            }

            return result
        except KeyError as e:
            raise ValueError(f"Missing expected field: {e}") from e

    def connectPostgres(self,db_info):
        
        username = db_info.get("username")
        password = quote_plus(db_info.get("password"))
        url = db_info.get("url")
        db_name = db_info.get("db")
        port = db_info.get("port")
        engine = db_info.get("engine").lower()  

        
        connection_url = f"postgresql://{username}:{password}@{url}:{port}/{db_name}"
        
        try:
            db_engine = create_engine(connection_url)
            print("Database engine created successfully.")
            return db_engine
        except Exception as e:
            print(e)
            raise RuntimeError("Failed to create SQLAlchemy engine.") from e

    def NewConnection(self,secret):

        # creds = fun(secret) genereate creds from google secret manager

        fields = self.parse_data_source_info(creds["secret_value"])
        dbMap = {}
        dbMap["POSTGRES"] = self.connectPostgres
        return dbMap[fields.get("engine")](fields)
    

