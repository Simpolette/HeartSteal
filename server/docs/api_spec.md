# API Specification

This document serves as the canonical reference for all backend API endpoints. Update this file whenever new endpoints are added or modified.

## Template

### [Endpoint Name]
-   **Method:** `[GET | POST | PUT | DELETE]`
-   **Route:** `[URL path]`
-   **Description:** `[Brief description]`
-   **Auth Required:** `[Yes | No]`

1.  **Request Parameters:**
    -   `[param_name]`: `[type]` - `[description]`

2.  **Request Body:**
    ```json
    {
      "key": "value",
      "optional_key": "value"
    }
    ```

3.  **Response (Success):**
    -   **Code:** `200 OK` / `201 Created`
    -   **Body:**
        ```json
        {
          "message": "Success message",
          "data": {
            "id": "...",
            "..."
          }
        }
        ```

4.  **Response (Error):**
    -   **Code:** `400 Bad Request` / `500 Internal Server Error`
    -   **Body:**
        ```json
        {
          "message": "Error description"
        }
        ```

---

## Existing Endpoints

### [Example] Register User
-   **Method:** `POST`
-   **Route:** `/api/v1/auth/signup`
-   **Description:** Registers a new user.
-   **Auth Required:** No

1.  **Request Body:**
    ```json
    {
      "username": "johndoe",
      "email": "john@example.com",
      "password": "strongPassword123"
    }
    ```

2.  **Response (Success):**
    -   **Code:** `201 Created`
    -   **Body:**
        ```json
        {
          "message": "User registered successfully"
        }
        ```
