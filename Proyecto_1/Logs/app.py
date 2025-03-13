from fastapi import FastAPI,  HTTPException, Body
from pydantic import BaseModel
import matplotlib.pyplot as plt
from datetime import datetime
from fastapi.responses import JSONResponse

app = FastAPI()

DATA = [
  {
    "time": "2023-10-21T12:00:00Z",
    "CPU": 14.11,
    "Mem": 0.65,
    "IO": 3463,
    "Disc": 0
  },
  {
    "time": "2023-10-21T12:00:10Z",
    "CPU": 0.45,
    "Mem": 0.69,
    "IO": 33,
    "Disc": 2000
  },  
  {
    "time": "2023-10-21T12:00:25Z",
    "CPU": 20.05,
    "Mem": 0.83,
    "IO": 8900,
    "Disc": 2120
  },
  {
    "time": "2023-10-21T12:00:35Z",
    "CPU": 13.05,
    "Mem": 0.0,
    "IO": 89,
    "Disc": 4120
  },
  {
    "time": "2023-10-21T12:00:45Z",
    "CPU": 1.05,
    "Mem": 0.80,
    "IO": 890,
    "Disc": 4122
  },
  {
    "time": "2023-10-21T12:00:55Z",
    "CPU": 12.05,
    "Mem": 0.80,
    "IO": 290,
    "Disc": 1122
  }
]



@app.get("/data")
def Get_Json():
    return JSONResponse(DATA)

