from fastapi import FastAPI, Request
from pydantic import BaseModel
from uuid import UUID
from datetime import datetime

class WebhookResponse(BaseModel):
    result: str

class Basicmessage(BaseModel):
    connection_id: UUID
    message_id: UUID
    content: str
    state: str
    sent_time: datetime

app = FastAPI()

@app.post("/topic/basicmessages", response_model=WebhookResponse, status_code=200)
async def handle_basic_message(message: Basicmessage):
    print(message)
    return {"result": "ok"}