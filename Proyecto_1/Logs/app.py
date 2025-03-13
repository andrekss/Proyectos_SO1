from fastapi import FastAPI,  HTTPException, Body
from pydantic import BaseModel
import matplotlib.pyplot as plt
from datetime import datetime
from fastapi.responses import JSONResponse

app = FastAPI()

DATA = [
  {
    "CPU": 0.18,
    "Disc": 0,
    "IO": 20,
    "Mem": 0.4,
    "time": "2025-03-13T09:33:51.002030121+00:00"
  },
  {
    "CPU": 3.42,
    "Disc": 0,
    "IO": 11292,
    "Mem": 0.0,
    "time": "2025-03-13T09:33:51.102113786+00:00"
  },
  {
    "CPU": 7.05,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.65,
    "time": "2025-03-13T09:33:51.202210735+00:00"
  },
  {
    "CPU": 6.98,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.05,
    "time": "2025-03-13T09:33:51.302330186+00:00"
  },
  {
    "CPU": 4.99,
    "Disc": 0,
    "IO": 13956,
    "Mem": 0.0,
    "time": "2025-03-13T09:33:51.402385138+00:00"
  },
  {
    "CPU": 3.09,
    "Disc": 0,
    "IO": 14591,
    "Mem": 0.0,
    "time": "2025-03-13T09:33:51.502478801+00:00"
  },
  {
    "CPU": 2.33,
    "Disc": 2037,
    "IO": 8440,
    "Mem": 0.64,
    "time": "2025-03-13T09:33:51.602578596+00:00"
  },
  {
    "CPU": 7.31,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.05,
    "time": "2025-03-13T09:33:51.702685623+00:00"
  },
  {
    "CPU": 6.99,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.0,
    "time": "2025-03-13T09:33:51.802784536+00:00"
  },
  {
    "CPU": 3.03,
    "Disc": 0,
    "IO": 16432,
    "Mem": 0.0,
    "time": "2025-03-13T09:33:51.902873441+00:00"
  },
  {
    "CPU": 6.96,
    "Disc": 0,
    "IO": 0,
    "Mem": 0.14,
    "time": "2025-03-13T09:33:52.002958859+00:00"
  }
]
#http://192.168.1.15:8000/data

@app.get("/data")
def Get_Json():
    return JSONResponse(DATA)

# docker run -d --name grafana -p 3000:3000 grafana/grafana:latest
