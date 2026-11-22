import os

import jwt
from fastapi import Depends, HTTPException
from fastapi.security import HTTPAuthorizationCredentials, HTTPBearer
from sqlalchemy import text
from sqlalchemy.engine import Connection

from app.db import get_connection

bearer = HTTPBearer(auto_error=False)
jwt_secret = os.environ["JWT_SECRET"]


def unauthorized() -> HTTPException:
    return HTTPException(
        status_code=401,
        detail="Invalid or expired access token",
        headers={"WWW-Authenticate": "Bearer"},
    )


def require_access_token(
    credentials: HTTPAuthorizationCredentials | None = Depends(bearer),
    connection: Connection = Depends(get_connection),
) -> str:
    if credentials is None:
        raise unauthorized()

    try:
        claims = jwt.decode(
            credentials.credentials,
            jwt_secret,
            algorithms=["HS256"],
            issuer="food-delivery-auth",
            audience="orders-api",
            options={"require": ["exp", "sub", "jti", "iss", "aud"]},
        )
    except jwt.InvalidTokenError:
        raise unauthorized() from None

    subject = claims["sub"]
    token_id = claims["jti"]
    if not subject.strip() or not token_id.strip():
        raise unauthorized()

    revoked = connection.execute(
        text("SELECT EXISTS (SELECT 1 FROM revoked_tokens WHERE jti = :jti)"),
        {"jti": token_id},
    ).scalar_one()

    if revoked:
        raise unauthorized()

    return subject