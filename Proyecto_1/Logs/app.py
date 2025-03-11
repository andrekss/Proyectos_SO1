from fastapi import FastAPI,  HTTPException, Body
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

@app.get("/fix")
def Recibir_Logs():
    return "hola mundo"

@app.post("/logs")
def Recibir_Logs(log: LogData):
    data = []
    if os.path.exists(LOG_FILE):
        with open(LOG_FILE, "r") as f:
            try:
                data = json.load(f)
            except:
                data = []

    data.append(log.dict()) # nuevo log

    # Guardar
    with open(LOG_FILE, "w") as f:
        json.dump(data, f, indent=2)
    #crear volumen afuera en referencia 
    return {"status": "ok", "message": "Log received"}

@app.post("/Generar_Gráfica")
def Generar_Gráficas():
    if not os.path.exists(LOG_FILE):
        return {"status": "error", "message": "No logs to generate graphs"}
    with open(LOG_FILE, "r") as f:
        data = json.load(f)

    return {"status": "ok", "message": "Graphs generated"}
