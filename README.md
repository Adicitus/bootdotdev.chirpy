# Chirpy
Chirpy is a Twitter/Bluesky-style micro-blogging platform developed for the Boot.dev course "Learn HTTP Servers".

This system as implemented in the course is primarily a RESTful HTTP API for for publishing and retrieving "Chirps".

## Installation

## Configuration

## API

- /admin: RESTful endpoints for application metadata and admin functions.
    - /metrics: Publishes metrics of the number of static assets served.
    - /reset: Resets the application, removing all data and resetting asset counters. *dev only*
- /api: RESTful endpoints for application functionality.
    - /chirps: Endpoint for Chirps.
        - GET: Retrieves all chirps in the system.
            - Query params:
                - order: Specify either "asc" (for oldest  first) or "desc" (for newest first). If not specified, uses "desc"
                - author_id: Specifies an ID (UUID) for a user whose chirps we want to view. If not specified, displays chirps for all users.
        - POST: Adds a new Chirp to the system.
            - Body: 
        - /{chirpID}: Operations on a Chirp specified by the given chirpID.
            - GET: Returns the specified Chirp with a status code of 200, or 404 if the ID doesn't match a Chirp
            - DELETE: Deletes the given Chirp and reponds with status 201, or 404 if the ID doesn't match a Chirp
    - /healthz: API health endpoint, returns status code 200 and body "OK" if the API is running.
    - /login: Authentication endpoint, used to request a new refresh token and an access token.
        - POST: Used to authenticate
            - Body:
            - 200: Login successful
                - Body:
            - 400:
                - Body:
            - 401: Login failed,
                - Body:
    - /refresh: Used to retrieve a new access token.
        - POST: Use a refresh token to retrieve a new access token.
    - /revoke: Used to revoke a refresh token, effectively loging out the user.
    - /validate_chirp: Used to validate a Chirp before submitting it.
        - POST: This endpoint only accepts POST requests
            - Body:
    - /users: Endpoint for Chirpy Users.
        - GET: Retrieves a list of users
        - POST: Used to create a new user
        - /{userID}: Performs operations on the user specified by userID (UUID).
            - GET: Retrieves the user's details
            - PUT: Update the user's details
            - DELTE: Removes the user
- /app: Serves static files for the application.
