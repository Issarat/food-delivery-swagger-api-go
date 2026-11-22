from datetime import date as Date
from typing import Literal

from fastapi import APIRouter, Depends
from pydantic import BaseModel
from sqlalchemy import text
from sqlalchemy.engine import Connection

from app.auth import require_access_token
from app.db import get_connection

router = APIRouter(
    prefix="/analytics/orders",
    tags=["Order Analytics"],
    dependencies=[Depends(require_access_token)],
)


class OrderBreakdownItem(BaseModel):
    metro_city: str
    order_channel: str
    cuisine_category: str
    total_orders: int


class DailyOrderCount(OrderBreakdownItem):
    date: Date


class MarketShareItem(BaseModel):
    value: str
    total_orders: int
    market_share_pct: float


def build_order_filter(
    *,
    metro_city: str | None,
    order_channel: str | None,
    cuisine_category: str | None,
) -> tuple[str, dict[str, str]]:
    filters = {
        "metro_city": metro_city,
        "order_channel": order_channel,
        "cuisine_category": cuisine_category,
    }
    params = {
        column: value.strip()
        for column, value in filters.items()
        if value is not None and value.strip()
    }

    if not params:
        return "", params

    conditions = [f"{column} = :{column}" for column in params]
    return "WHERE " + " AND ".join(conditions), params


@router.get("/summary", response_model=list[OrderBreakdownItem])
@router.get("/breakdown", response_model=list[OrderBreakdownItem])
def get_order_breakdown(
    connection: Connection = Depends(get_connection),
    metro_city: str | None = None,
    order_channel: str | None = None,
    cuisine_category: str | None = None,
) -> list[OrderBreakdownItem]:
    where_clause, params = build_order_filter(
        metro_city=metro_city,
        order_channel=order_channel,
        cuisine_category=cuisine_category,
    )

    query = text(f"""
        SELECT metro_city,
               order_channel,
               cuisine_category,
               COUNT(*) AS total_orders
        FROM orders
        {where_clause}
        GROUP BY metro_city, order_channel, cuisine_category
        ORDER BY total_orders DESC, metro_city, order_channel, cuisine_category
    """)

    rows = connection.execute(query, params).mappings()
    return [OrderBreakdownItem(**row) for row in rows]


@router.get("/daily", response_model=list[DailyOrderCount])
def get_daily_order_counts(
    connection: Connection = Depends(get_connection),
    metro_city: str | None = None,
    order_channel: str | None = None,
    cuisine_category: str | None = None,
) -> list[DailyOrderCount]:
    where_clause, params = build_order_filter(
        metro_city=metro_city,
        order_channel=order_channel,
        cuisine_category=cuisine_category,
    )

    query = text(f"""
        SELECT (created_at AT TIME ZONE 'Asia/Bangkok')::date AS date,
               metro_city,
               order_channel,
               cuisine_category,
               COUNT(*) AS total_orders
        FROM orders
        {where_clause}
        GROUP BY date, metro_city, order_channel, cuisine_category
        ORDER BY date, total_orders DESC,
                 metro_city, order_channel, cuisine_category
    """)

    rows = connection.execute(query, params).mappings()
    return [DailyOrderCount(**row) for row in rows]


@router.get("/marketshare", response_model=list[MarketShareItem])
def get_order_market_share(
    group_by: Literal["metro_city", "order_channel", "cuisine_category"],
    connection: Connection = Depends(get_connection),
    metro_city: str | None = None,
    order_channel: str | None = None,
    cuisine_category: str | None = None,
) -> list[MarketShareItem]:
    where_clause, params = build_order_filter(
        metro_city=metro_city,
        order_channel=order_channel,
        cuisine_category=cuisine_category,
    )

    query = text(f"""
        SELECT {group_by} AS value,
               COUNT(*) AS total_orders,
               ROUND(100.0 * COUNT(*) / SUM(COUNT(*)) OVER (), 2)
                   AS market_share_pct
        FROM orders
        {where_clause}
        GROUP BY {group_by}
        ORDER BY total_orders DESC, value ASC
    """)

    rows = connection.execute(query, params).mappings()
    return [MarketShareItem(**row) for row in rows]