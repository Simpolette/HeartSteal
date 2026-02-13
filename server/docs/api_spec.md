# API Specification

This document serves as the canonical reference for all backend API endpoints. Update this file whenever new endpoints are added or modified.

---

## Auth Endpoints

### Register User
-   **Method:** `POST`
-   **Route:** `/api/signup`
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

3.  **Response (Error):**
    -   **Code:** `409 Conflict` — Email already exists
    -   **Code:** `400 Bad Request` — Validation error
    -   **Body:**
        ```json
        {
          "message": "Error description"
        }
        ```

---

### Login
-   **Method:** `POST`
-   **Route:** `/api/login`
-   **Description:** Authenticates a user and returns access + refresh tokens.
-   **Auth Required:** No

1.  **Request Body:**
    ```json
    {
      "username": "johndoe",
      "password": "strongPassword123"
    }
    ```

2.  **Response (Success):**
    -   **Code:** `200 OK`
    -   **Body:**
        ```json
        {
          "message": "Login successfully",
          "data": {
            "access_token": "eyJhbGciOi...",
            "refresh_token": "eyJhbGciOi..."
          }
        }
        ```

3.  **Response (Error):**
    -   **Code:** `401 Unauthorized` — Invalid credentials
    -   **Body:**
        ```json
        {
          "message": "Invalid email or password"
        }
        ```

---

### Forgot Password
-   **Method:** `POST`
-   **Route:** `/api/forgot-password`
-   **Description:** Sends a 6-digit PIN to the user's email for password recovery.
-   **Auth Required:** No

1.  **Request Body:**
    ```json
    {
      "email": "john@example.com"
    }
    ```

2.  **Response (Success):**
    -   **Code:** `200 OK`
    -   **Body:**
        ```json
        {
          "message": "If the email exists, a PIN code has been sent"
        }
        ```

> [!NOTE]
> Returns 200 even if the email does not exist to prevent email enumeration attacks.

---

### Verify PIN
-   **Method:** `POST`
-   **Route:** `/api/verify-pin`
-   **Description:** Verifies the 6-digit PIN and returns a temporary reset token.
-   **Auth Required:** No

1.  **Request Body:**
    ```json
    {
      "email": "john@example.com",
      "pin_code": "123456"
    }
    ```

2.  **Response (Success):**
    -   **Code:** `200 OK`
    -   **Body:**
        ```json
        {
          "message": "PIN verified successfully",
          "data": {
            "reset_token": "eyJhbGciOi..."
          }
        }
        ```

3.  **Response (Error):**
    -   **Code:** `401 Unauthorized` — Invalid or expired PIN
    -   **Body:**
        ```json
        {
          "message": "invalid PIN code"
        }
        ```

---

### Reset Password
-   **Method:** `POST`
-   **Route:** `/api/reset-password`
-   **Description:** Resets the user's password using the reset token from Verify PIN.
-   **Auth Required:** No (uses reset token)

1.  **Request Body:**
    ```json
    {
      "reset_token": "eyJhbGciOi...",
      "new_password": "newStrongPassword123"
    }
    ```

2.  **Response (Success):**
    -   **Code:** `200 OK`
    -   **Body:**
        ```json
        {
          "message": "Password reset successfully"
        }
        ```

3.  **Response (Error):**
    -   **Code:** `401 Unauthorized` — Invalid reset token
    -   **Body:**
        ```json
        {
          "message": "invalid reset token"
        }
        ```

---

### Sign Out
-   **Method:** `POST`
-   **Route:** `/api/signout`
-   **Description:** Signs out the user by invalidating all session tokens.
-   **Auth Required:** Yes (Bearer token)

1.  **Request Headers:**
    ```
    Authorization: Bearer <access_token>
    ```

2.  **Response (Success):**
    -   **Code:** `200 OK`
    -   **Body:**
        ```json
        {
          "message": "Signed out successfully"
        }
        ```

3.  **Response (Error):**
    -   **Code:** `401 Unauthorized` — Not authorized
