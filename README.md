# imagestore
Simple Rest API for storing images, written in Go.

## Setup
Setup is simple, just run:
```
git clone https://github.com/unitezen/imagestore.git
cd imagestore
sudo docker compose up -d
```
and the API will be available on the host machine at http://localhost:8080

## Features
- Fiber (HTTP framework) + Grom (ORM) with PostgreSQL database
- Custom error handler for unified JSON error message with status code
- JSON request payload validation using go-playground/validator
- Authentication via middleware using X-API-Key header for private endpoints
- Hashed password storage with BCrypt
- Image storage as database entry using Base64
- Verification of uploaded base64 data to target image/jpeg and image/png content type

## Endpoints
| Endpoint     | Method | Request Payload                                                                          | Response          | Authentication Required |
|--------------|--------|------------------------------------------------------------------------------------------|-------------------|-------------------------|
| /users       | POST   | username (min=4, max=32, alphanumeric), password (min=8, max=32), email (max=255, email) | api_key           | No                      |
| /sessions    | POST   | username (min=4, max=32, alphanumeric), password (min=8, max=32)                         | api_key           | No                      |
| /sessions    | DELETE |                                                                                          | success           | Yes                     |
| /images      | GET    |                                                                                          | images (id, name) | Yes                     |
| /images      | POST   | name (max=255), data (base64 JPEG/PNG image)                                             | id                | Yes                     |
| /images/{id} | GET    |                                                                                          | name, data        | Yes                     |
| /images/{id} | PATCH  | name (max=255)                                                                           | id, name          | Yes                     |
| /images/{id} | DELETE |                                                                                          | success           | Yes                     |

- Create new user with `POST /users`
- Login with `POST /sessions`
- Logout with `DELETE /sessions`
- List all images with `GET /images` and individual image with `GET /images/{id}`
- Upload a base64-encoded JPEG/PNG image with `POST /images`
- Modify image name with `PATCH /images/{id}`
- Delete image with `DELETE /images/{id}`

## Possible Improvements
- Pass in configuration parameters via environment variables (database connection string is currently hardcoded)
- More robust error handling and error messages
- Add more code comment and docstrings
- Refactor to reduce repetitive parsing/error checking code
- Improve project structure
- Unit and end-to-end tests
