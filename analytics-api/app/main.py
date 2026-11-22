from fastapi import Depends, FastAPI
from sqlalchemy import text
from sqlalchemy.engine import Connection

from app.db import get_connection
from app.router import router as analytics_router

app = FastAPI(
    title="Food Delivery Analytics API",
    description=(
        "Get an access token from the "
        "[Go API Swagger](http://localhost:8080/swagger/), "
        "then paste it into **Authorize** on this page."
    ),
)

@app.get("/health/db")
def check_database(
    connection: Connection = Depends(get_connection),
) -> dict[str, int]:
    result = connection.execute(text("SELECT 1")).scalar_one()
    return {"database_check": result}

app.include_router(analytics_router)