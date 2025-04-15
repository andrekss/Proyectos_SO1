from locust import HttpUser, task, between
import random

# Posibles países y tipos de clima
Paises = ["GT", "US", "MX", "BR", "FR", "DE", "JP", "CA", "AR", "PE"]
Climas = [ "Nubloso", "Soleado", "Lluvioso"]
Descripciones = ["Está lloviendo", "El cielo está despejado", "Hay neblina", "Está nublado", "Hace calor"]

class Usuario_Clima(HttpUser):
    tiempo_espera = between(0.1, 0.3)  # Tiempo entre peticiones (ajustable)

    @task
    def Enviar_Tweet(self):
        payload = {
            "Description": random.choice(Descripciones),
            "Country": random.choice(Paises),
            "Weather": random.choice(Climas)
        }

        self.client.post("/input", json=payload) # Api expuesta
