## Overview

This service is a high-performance identity and access management component to authenticate users. 
It handles the complete authentication lifecycle, ensuring secure communication between services through gRPC
and managing user sessions with low-latency caching.

## Key Features

1. Security: Password hashing and secure storage of user credentials.
2. Token Management: Issuing and validating JWT (Access/Refresh) tokens.
3. Centralized Authentication: Acting as a source of truth for user identity across the platform.

## ENV Variables

#### db
`DB_NAME` - The name of the database to connnect to.

#### jwt
`JWT_SECRET` - The secret key used to sign JWT tokens.\
`JWT_ACCESS_TOKEN_TTL` - The TTL for access tokens.\
`JWT_REFRESH_TOKEN_TTL` - The TTL for refresh tokens.\
`JWT_ISSUER` - The issuer of the JWT tokens.

## Endpoints

[//]: # (TODO: Add endpoints)