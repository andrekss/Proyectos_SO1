'''
Number of users: 10
Ramp up: 1
Host: https://34.58.7.218.nip.io
'''

from locust.contrib.fasthttp import FastHttpUser
from locust import task, between
import random

# Opciones para el payload
Paises = ["GT", "US", "MX", "BR", "FR", "DE", "JP", "CA", "AR", "PE"]
Climas = ["Nubloso", "Soleado", "Lluvioso"]
Descripciones = ["Está lloviendo", "El cielo está despejado", "Hay neblina", "Está nublado", "Hace calor"]

class Usuario_Clima(FastHttpUser):
    wait_time =  lambda self: 0
    last_payload = None
    intentos = 0
    LIMITE = 10000

    @task
    def enviar_tweet(self):
        if self.intentos < self.LIMITE:
            nuevo_payload = self.generar_payload_diferente()
            if nuevo_payload:
                self.client.post("/input", json=nuevo_payload, verify=False)
                self.last_payload = nuevo_payload
                self.intentos += 1
        else:
            self.environment.runner.quit()

    def generar_payload_diferente(self):
        while True:
            payload = {
                "description": random.choice(Descripciones),
                "country": random.choice(Paises),
                "weather": random.choice(Climas)
            }
            if payload != self.last_payload:
                return payload
