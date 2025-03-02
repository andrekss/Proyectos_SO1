from fastapi import FastAPI
from pydantic import BaseModel
import json
import matplotlib.pyplot as plt
from datetime import datetime
import os

app = FastAPI()

LOG_FILE = "/logs/logs.json"  # Volumen montado

class LogData(BaseModel):
    container_id: str
    timestamp: str
    event: str
    # y los campos que necesites

@app.post("/logs")
def Recibir_Logs(log: LogData):
    # 1. Leer archivo JSON actual
    data = []
    if os.path.exists(LOG_FILE):
        with open(LOG_FILE, "r") as f:
            try:
                data = json.load(f)
            except:
                data = []

    # 2. Agregar el nuevo log
    data.append(log.dict())

    # 3. Guardar
    with open(LOG_FILE, "w") as f:
        json.dump(data, f, indent=2)

    return {"status": "ok", "message": "Log received"}

@app.post("/generate-graphs")
def Generar_Gráficas():
    # 1. Leer datos
    if not os.path.exists(LOG_FILE):
        return {"status": "error", "message": "No logs to generate graphs"}
    with open(LOG_FILE, "r") as f:
        data = json.load(f)

    return {"status": "ok", "message": "Graphs generated"}
