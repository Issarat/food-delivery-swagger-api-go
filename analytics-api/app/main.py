from fastapi import Depends, FastAPI
from sqlalchemy import text
from sqlalchemy.engine import Connection

from app.db import get_connection

app = FastAPI()


@app.get("/health/db")
def check_database(
    connection: Connection = Depends(get_connection),
) -> dict[str, int]:
    result = connection.execute(text("SELECT 1")).scalar_one()
    return {"database_check": result}