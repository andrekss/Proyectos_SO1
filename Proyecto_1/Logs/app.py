from fastapi import FastAPI,  HTTPException, Body
from pydantic import BaseModel
import matplotlib.pyplot as plt
from datetime import datetime
from fastapi.responses import JSONResponse

app = FastAPI()

class DataItem(BaseModel):
    CPU: float
    Disc: int
    IO: int
    Mem: float
    time: str  

DATA = []

# referenicas
"""[
  {
    "CPU": 0.52,
    "Disc": 220,
    "IO": 3238,
    "Mem": 2.09,
    "time": "2025-03-13T20:59:44.275321350+00:00"
  },
  {
    "CPU": 0.13,
    "Disc": 0,
    "IO": 25,
    "Mem": 0.4,
    "time": "2025-03-13T20:59:44.285911290+00:00"
  },
  {
    "CPU": 0.36,
    "Disc": 220,
    "IO": 3255,
    "Mem": 2.1,
    "time": "2025-03-13T20:59:55.731589559+00:00"
  },
  {
    "CPU": 0.11,
    "Disc": 0,
    "IO": 45,
    "Mem": 0.41,
    "time": "2025-03-13T20:59:55.741678699+00:00"
  },
  {
    "CPU": 0.39,
    "Disc": 220,
    "IO": 3256,
    "Mem": 2.12,
    "time": "2025-03-13T21:00:50.322695750+00:00"
  },
  {
    "CPU": 0.09,
    "Disc": 0,
    "IO": 45,
    "Mem": 0.41,
    "time": "2025-03-13T21:00:50.332769362+00:00"
  },
  {
    "CPU": 7.01,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.0,
    "time": "2025-03-13T21:00:50.342838496+00:00"
  },
  {
    "CPU": 7.15,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.0,
    "time": "2025-03-13T21:00:50.352921966+00:00"
  },
  {
    "CPU": 6.97,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.65,
    "time": "2025-03-13T21:00:50.363017718+00:00"
  },
  {
    "CPU": 4.27,
    "Disc": 1468,
    "IO": 4223,
    "Mem": 0.84,
    "time": "2025-03-13T21:00:50.373105877+00:00"
  },
  {
    "CPU": 6.93,
    "Disc": 1198,
    "IO": 3755,
    "Mem": 0.83,
    "time": "2025-03-13T21:00:50.383174430+00:00"
  },
  {
    "CPU": 6.95,
    "Disc": 1023,
    "IO": 2942,
    "Mem": 0.69,
    "time": "2025-03-13T21:00:50.393249915+00:00"
  },
  {
    "CPU": 0.06,
    "Disc": 0,
    "IO": 116,
    "Mem": 0.0,
    "time": "2025-03-13T21:00:50.403341059+00:00"
  },
  {
    "CPU": 0.36,
    "Disc": 0,
    "IO": 92,
    "Mem": 0.0,
    "time": "2025-03-13T21:00:50.413426022+00:00"
  },
  {
    "CPU": 6.97,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.65,
    "time": "2025-03-13T21:00:50.423538977+00:00"
  },
  {
    "CPU": 7.03,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.0,
    "time": "2025-03-13T21:00:50.433623810+00:00"
  },
  {
    "CPU": 0.46,
    "Disc": 221,
    "IO": 3280,
    "Mem": 2.19,
    "time": "2025-03-13T21:01:32.754959194+00:00"
  },
  {
    "CPU": 0.1,
    "Disc": 0,
    "IO": 45,
    "Mem": 0.41,
    "time": "2025-03-13T21:01:32.765041812+00:00"
  },
  {
    "CPU": 7.03,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.0,
    "time": "2025-03-13T21:01:32.775132606+00:00"
  },
  {
    "CPU": 6.98,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.0,
    "time": "2025-03-13T21:01:32.785222548+00:00"
  },
  {
    "CPU": 6.77,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.3,
    "time": "2025-03-13T21:01:32.795309344+00:00"
  },
  {
    "CPU": 4.94,
    "Disc": 0,
    "IO": 37726,
    "Mem": 0.0,
    "time": "2025-03-13T21:01:32.805434662+00:00"
  },
  {
    "CPU": 7.36,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.0,
    "time": "2025-03-13T21:01:32.815595616+00:00"
  },
  {
    "CPU": 6.92,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.0,
    "time": "2025-03-13T21:01:32.825692270+00:00"
  },
  {
    "CPU": 6.99,
    "Disc": 3231,
    "IO": 13688,
    "Mem": 0.83,
    "time": "2025-03-13T21:01:32.835777053+00:00"
  },
  {
    "CPU": 6.82,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.65,
    "time": "2025-03-13T21:01:32.845861935+00:00"
  },
  {
    "CPU": 6.94,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.0,
    "time": "2025-03-13T21:01:32.855963529+00:00"
  },
  {
    "CPU": 6.99,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.0,
    "time": "2025-03-13T21:01:32.866038133+00:00"
  },
  {
    "CPU": 1.12,
    "Disc": 0,
    "IO": 223,
    "Mem": 0.0,
    "time": "2025-03-13T21:01:32.876129057+00:00"
  },
  {
    "CPU": 6.78,
    "Disc": 37,
    "IO": 379,
    "Mem": 0.83,
    "time": "2025-03-13T21:01:32.886225160+00:00"
  },
  {
    "CPU": 0.34,
    "Disc": 0,
    "IO": 6,
    "Mem": 0.0,
    "time": "2025-03-13T21:01:32.896308961+00:00"
  }
]"""


#http://192.168.1.15:8000/data
#http://192.168.1.5:8000/data # para laptop

@app.post("/Cargar_Json")
def cargar_json(data: list[DataItem]):
    global DATA
    # Convierte los objetos Pydantic a diccionarios, si es necesario.
    DATA = [item.dict() for item in data]
    return "Datos actualizados en tiempo real"

@app.get("/data")
def Get_Json():
    return JSONResponse(DATA)

# docker run -d --name grafana -p 3000:3000 grafana/grafana:latest
