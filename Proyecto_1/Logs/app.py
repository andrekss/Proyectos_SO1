from fastapi import FastAPI,  HTTPException
from pydantic import BaseModel
import json
import matplotlib.pyplot as plt
from datetime import datetime
import subprocess
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


@app.post("/setCronJob")
def set_cronjob(req: int):
    if req == 1:
        cron_content = (
            "* * * * * /home/andres/Proyectos_SO1/Sistemas_Operativos_1/Proyectos_SO1/Proyecto_1/Scripts/Contenedores.sh\n"
            "* * * * * sleep 30 && /home/andres/Proyectos_SO1/Sistemas_Operativos_1/Proyectos_SO1/Proyecto_1/Scripts/Contenedores.sh\n"
        )
        cron_file = "/tmp/mycron"
        with open(cron_file, "w") as f:
            f.write(cron_content)
        subprocess.run(["crontab", cron_file], check=True)
        return {"status": "ok", "message": "Crontab actualizado con el script cada 30 segundos"}
    elif req == 0:
        subprocess.run(["crontab", "-r"], check=True)
        return {"status": "ok", "message": "Crontab borrado"}
    else:
        raise HTTPException(status_code=400, detail="Valor de acción inválido. Use 1 para crear y 0 para borrar.")

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

    return {"status": "ok", "message": "Log received"}

@app.post("/Generar_Gráfica")
def Generar_Gráficas():
    if not os.path.exists(LOG_FILE):
        return {"status": "error", "message": "No logs to generate graphs"}
    with open(LOG_FILE, "r") as f:
        data = json.load(f)

    return {"status": "ok", "message": "Graphs generated"}
