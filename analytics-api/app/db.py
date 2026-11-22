import os
from collections.abc import Iterator

from dotenv import load_dotenv
from sqlalchemy import create_engine
from sqlalchemy.engine import Connection, URL

load_dotenv()

database_url = os.getenv("DATABASE_URL")
if not database_url:
    database_url = URL.create(
        drivername="postgresql+psycopg",
        username=os.environ["POSTGRES_USER"],
        password=os.environ["POSTGRES_PASSWORD"],
        host=os.environ["POSTGRES_HOST"],
        port=int(os.getenv("POSTGRES_PORT", "5432")),
        database=os.environ["POSTGRES_DB"],
    )

engine = create_engine(database_url, pool_pre_ping=True)

def get_connection() -> Iterator[Connection]:
    with engine.connect() as connection:
        yield connection