# Chirpy
Chirpy is a Twitter/Bluesky-style micro-blogging platform developed for the Boot.dev course "Learn HTTP Servers".

This system as implemented in the course is primarily a RESTful HTTP API for for publishing and retrieving "Chirps".

## Installation
Chirpy is designed to use [https://www.postgresql.org/](PostgreSQL) as a backend, so you will need to install it before running Chirpy.

To install Chirpy, download the project files:
```
> git clone https://github.com/Adicitus/bootdotdev.chirpy.git
```

## Configuration
In order to run the chirpy server, you'll need a .env file.

For details on configuration options, see .env.sample in the source files.

## Basic Usage
Simply Launch the server from inside the source directory:
```
go run .
```

## API
Chirpy is primarily a RESTful API, serving the following endpoints:

- **/admin**: RESTful endpoints for application metadata and admin functions.
    - **/metrics**: Publishes metrics of the number of static assets served.
        - **GET**: returns a HTML document with the number of static assets that have been served during this servers lifetime.
            - Security: Requires an access token.
            - response:
                - 200: Returns the metrics report.
    - **/reset**: Allows admins to reset the application state.
        - **POST**: Resets the application, removing all data and resetting asset counters. *dev mode only*
            - Response:
                - 200: Reset successfull
                - 403: Server is not in dev mode
- **/api**: RESTful endpoints for application functionality.
    - **/chirps**: Endpoint for Chirps.
        - **GET**: Retrieves all chirps in the system.
            - Security: None
            - Query params:
                - order: Specify either "asc" (for oldest  first) or "desc" (for newest first). If not specified, uses "desc"
                - author_id: Specifies an ID (UUID) for a user whose chirps we want to view. If not specified, displays chirps for all users.
            - Response:
                - 200:
        - **POST**: Adds a new Chirp to the system.
            - Body: *application/json*
            ```
            {
                "body": "<Text for the chirp, max 140 characters, minimum 1>"
            }
            ```
            - Response:
                - 201: Chirp created successfully
                    - Body: *application/json*
                    ```
                    {
                        "id": "<UUID of the newly created chirp>",
                        "created_at": "<date and time that the chirp was created>",
                        "updated_at": "<date and time that the chirp was last updated>",
                        "body": "<The submitted body>",
                        "user_id": "<The id of the user that created the chirp>"
                    }
                    ```
                - 400: Chirp violates the requirements (too long or too short).
        - **/{chirpID}**: Operations on a Chirp specified by the given chirpID.
            - **GET**: Retrieves the specified chirp
                - Security: None
                - Response:
                    - 200: Chirp found
                        - Body: application/json
                        ```
                        {
                            "id": "<UUID of the newly created chirp>",
                            "created_at": "<date and time that the chirp was created>",
                            "updated_at": "<date and time that the chirp was last updated>",
                            "body": "<The submitted body>",
                            "user_id": "<The id of the user that created the chirp>"
                        }
                        ```
                    - 404: The specified chirp does not exist.
            - **DELETE**: Deletes the given Chirp and reponds with status 201, or 404 if the ID doesn't match a Chirp
                - Security: Requires an access token
                - Response:
                    - 204: Chirp removed
                    - 403: The calling user is not the owner of the chirp
                    - 404: The specified chirp does not exist
    - **/healthz**: API health endpoint, returns status code 200 and body "OK" if the API is running.
        - Security: None
        - Response:
            - 200: Server is alive.
                - Body: text/plain
                ```
                OK
                ```
    - **/login**: Authentication endpoint, used to request a new refresh token and an access token.
        - **POST**: Used to authenticate
            -Security: None
            - Body: application/json
              ```
              {
                "email": "<User's registered email>",
                "password": "<User's password>"
              }
              ```
            - Response:
                - 200: Login successful
                    - Body: application/json
                    ```
                    {
                        "id": "<User ID>",
                        "created_at": "<Date that the user profile was created>",
                        "updated_at": "<Date that the user profile was last updated>",
                        "email": "<User's registered email>",
                        "token": "<A new access token, valid for 1 hour>",
                        "refresh_roken": "<A new refresh token, valid for 60 days>"
                    }
                    ```
                - 400: Invalid details supplied.
                - 401: Login failed.
    - **/refresh**: Used to retrieve a new access token.
        - **POST**: Use a refresh token to retrieve a new access token.
            - Security: Requires a refresh token
            - Response:
                - 200: New access token issued
                    - Body: application/json
                    ```
                    {
                        "token": "<The newly created access token>"
                    }
                    ```
    - **/revoke**: Used to revoke a refresh token, effectively loging out the user.
        - Security: Requires a refresh token.
        - Response:
            - 204: Refresh token revoked.
    - **/validate_chirp**: Used to validate a Chirp before submitting it.
        - POST: This endpoint only accepts POST requests
            - Security: Requires an access token
            - Body: application/json
            {
                "body": "<Text to validate>"
                "user_id":  "<ID of the user submitting the chirp>"
            }
            - Respone:
                - 200: Chirp is valid
                    - Body: application/json
                    ```
                    {
                        "cleaned_body": "<The submitted body with forbidden words censored>",
                        "user_id": "<ID of the user that submitted the chirp>"
                    }
                    ```
                - 400: Chirp could not be read
                - 409: Chirp is invalid
                    - Body: applicaion/json
                    {
                        "error": "<Reason for the invalidation>"
                    }
    - **/users**: Endpoint for Chirpy Users.
        - POST: Used to create a new user
        - **/{userID}**: Performs operations on the user specified by userID (UUID).
            - PUT: Update the user's details
                - Security: Requires access token
                - Body: applicaiton/json
                ```
                {
                    "email": "<New email for the user, can be omitted to retain old email>",
                    "password: "<New password for the user, can be omitted to retain old password>"
                }
                ```
                - Response:
                    - 200: 
                    - 404:
            - DELTE: Removes the user
- **/app**: Serves static files for the application.

### Notes on Security
- Any call that violates the security specification of an endpoint will result in a 401 status code.
- All user tokens are "Bearer" tokens.
- Webhooks use "ApiKey" tokens.
- If the server is in dev mode, the server will return error details for calls that result in non-OK status codes. Otherwise a generic status message will be returned instead.