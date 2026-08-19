````md
# Backend Integration Documentation

## Base URL

```txt
https://chat-app-ta.duckdns.org
```
````

---

# API Base URL

```txt
https://chat-app-ta.duckdns.org/api/v1
```

---

# Authentication System Overview

Backend does NOT use JWT Bearer Token on frontend.

Authentication uses:

- HTTP-only cookies
- Secure cookies
- Session-based authentication

Cookies:

- `session_id`
- `device_id`

Example cookie values:

```txt
session_id=123e4567-e89b-12d3-a456-426614174000
device_id=123e4567-e89b-12d3-a456-426614174000
```

Frontend requirements:

- Always use `withCredentials: true`
- Browser automatically sends cookies
- Do NOT manually store tokens in localStorage/sessionStorage

Recommended axios setup:

```ts
import axios from "axios";

export const api = axios.create({
  baseURL: "https://chat-app-ta.duckdns.org/api/v1",
  withCredentials: true,
});
```

---

# Auth Flow

## Login Flow

1. Frontend sends request to:

```http
POST /api/v1/auth/google/login
```

2. Backend communicates with Auth Service.

3. Backend sets:

- `session_id`
- `device_id`

as:

- HTTP-only cookies
- Secure cookies

4. Session data stored in:

- Redis
- Database

5. Next authenticated requests:

- Middleware reads cookies
- Validates session
- Injects:
  - `user_id`
  - `user_uuid`

into request context.

---

# API Security

Public API routes additionally require:

```http
X-Api-Key: <api_key>
```

Need confirmation from backend:

- which endpoints require `X-Api-Key`
- actual API key value for frontend environment

---

# Available APIs

# 1. Authentication APIs

## Login With Google

```http
POST /api/v1/auth/google/login
```

Full URL:

```txt
https://chat-app-ta.duckdns.org/api/v1/auth/google/login
```

Content-Type:

```txt
application/json
```

Request Body:

```json
{
  "auth_code": "string"
}
```

Purpose:

- Login using Google OAuth

Success Response:

```json
{
  "success": true,
  "data": {}
}
```

Possible Response:

```json
{
  "success": true,
  "message": "Login successful"
}
```

Error Response:

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": "Optional details"
}
```

Authentication Result:

- Backend automatically sets cookies:
  - `session_id`
  - `device_id`

---

## Logout Current Device

```http
DELETE /api/v1/auth/logout
```

Purpose:

- Logout current session/device

Auth:

- Requires cookies

---

## Logout All Devices

```http
DELETE /api/v1/auth/logout/all
```

Purpose:

- Logout all sessions/devices

Auth:

- Requires cookies

---

# 2. User APIs

## Get Current User Profile

```http
GET /api/v1/user/profile
```

Purpose:

- Get authenticated user profile

Auth:

- Requires cookies

Success Response:

```json
{
  "success": true,
  "data": {}
}
```

---

## Get User Profile By UUID

```http
GET /api/v1/user/profile/:uuid
```

Example:

```http
GET /api/v1/user/profile/123
```

Purpose:

- Get another user's profile

Auth:

- Requires cookies

---

## Update User Profile

```http
PUT /api/v1/user/profile
```

Content-Type:

```txt
multipart/form-data
```

Purpose:

- Update profile
- Upload avatar

Auth:

- Requires cookies

Multipart Form Fields:

| Field      | Type   | Required |
| ---------- | ------ | -------- |
| avatar_url | File   | No       |
| name       | string | No       |
| birthday   | string | No       |
| phone      | string | No       |

Avatar Upload:

- File uploads directly to Google Cloud Storage
- Bucket:
  - `GOOGLE_CLOUD_STORAGE_BUCKET_RAW`

Example FormData:

```ts
const formData = new FormData();

formData.append("avatar_url", file);
formData.append("name", "John Doe");
formData.append("birthday", "2000-01-01");
formData.append("phone", "0123456789");
```

---

## Disable User

```http
DELETE /api/v1/user/disable
```

Purpose:

- Disable/delete user account

Auth:

- Requires cookies

---

# 3. Avatar Report API

## Report Avatar

```http
POST /api/v1/user/report-avatar
```

Purpose:

- Report processed avatar

Auth:

- Requires cookies

Request Body:

```json
{
  "uuid": "string",
  "name": "string"
}
```

---

# 4. Notify / SSE APIs

## SSE Endpoint

```http
GET /api/v1/notify/sse
```

Purpose:

- Realtime notifications/events

Technology:

- Server-Sent Events (SSE)

Authentication:

- Uses same cookie authentication:
  - `session_id`
  - `device_id`

No JWT required.

Frontend Example:

```ts
const eventSource = new EventSource(
  "https://chat-app-ta.duckdns.org/api/v1/notify/sse",
  {
    withCredentials: true,
  },
);
```

---

# SSE Event Format

## Connection Success

```txt
: connected
```

---

## Ping / Heartbeat

Sent every 5 seconds:

```txt
: ping
```

---

## Data Event

Raw SSE Format:

```txt
data: {"user_id":"...","status":"...","file_path":"..."}
```

Parsed JSON Payload:

```json
{
  "user_id": "string",
  "status": "string",
  "file_path": "string"
}
```

Frontend Example:

```ts
eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data);

  console.log(data.user_id);
  console.log(data.status);
  console.log(data.file_path);
};
```

---

# Verify ID Token OTP

## VerifyIDTokenOTP

Content-Type:

```txt
application/json
```

Request Body:

```json
{
  "id_token": "string"
}
```

---

# Standard Response Format

# Success Response

## Success With Data

```json
{
  "success": true,
  "data": {}
}
```

---

## Success With Message

```json
{
  "success": true,
  "message": "Success message"
}
```

---

# Error Response

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": "Optional details"
}
```

---

# Validation Error Format

## HTTP Layer Validation Error

```json
{
  "error": {
    "field_name": "Validation message",
    "another_field": "Validation message"
  }
}
```

Example:

```json
{
  "error": {
    "name": "name là bắt buộc",
    "phone": "phone không hợp lệ"
  }
}
```

---

## gRPC Validation Error

```json
{
  "error": "VALIDATION_ERROR",
  "detail": {
    "field_name": "Error description"
  }
}
```

---

# Status Code Convention

| Error Code            | HTTP Status |
| --------------------- | ----------- |
| BAD_REQUEST           | 400         |
| UNAUTHORIZED          | 401         |
| NOT_FOUND             | 404         |
| CONFLICT              | 409         |
| TOO_MANY_REQUESTS     | 429         |
| INTERNAL_SERVER_ERROR | 500         |

Additional Notes:

- gRPC errors are mapped to HTTP status codes.
- Examples:
  - `InvalidArgument` -> 400
  - `PermissionDenied` -> 403
  - `Unauthenticated` -> 401

---

# Current Backend Modules

Currently available modules:

- AuthModule
- UserModule
- NotifyModule

---

# Missing Features / APIs

Currently NOT available:

- Message APIs
- Send message
- Get conversations
- Chat rooms
- Typing status
- Read receipts
- WebSocket APIs

Realtime currently only uses SSE.

---

# Backend Routes Summary

```txt
GET     /
GET     /healthz

POST    /api/v1/auth/google/login
DELETE  /api/v1/auth/logout
DELETE  /api/v1/auth/logout/all

GET     /api/v1/user/profile
GET     /api/v1/user/profile/:uuid
PUT     /api/v1/user/profile
DELETE  /api/v1/user/disable

POST    /api/v1/user/report-avatar

GET     /api/v1/notify/sse
```

---

# Recommended Frontend Structure

```txt
src/
├── api/
│   ├── auth.api.ts
│   ├── user.api.ts
│   └── notify.api.ts
│
├── services/
│   └── sse.service.ts
│
├── lib/
│   └── axios.ts
│
├── hooks/
│   └── useSSE.ts
│
├── context/
│   └── AuthContext.tsx
```

---

# Important Frontend Notes

## Cookies

Frontend MUST:

- enable credentials
- never manually manage session cookies
- rely on browser cookie handling

---

## SSE

Frontend MUST:

- reconnect on disconnect
- handle heartbeat events
- parse JSON safely

---

## File Upload

Frontend MUST:

- use `multipart/form-data`
- use `FormData`
- never send avatar as JSON

---

## Error Handling

Frontend SHOULD:

- support both validation error formats
- handle custom error codes
- map HTTP status properly

---

# Recommended Frontend Integration Priority

1. Google Login
2. Auth persistence
3. Current user profile
4. Protected routes
5. Update profile + avatar upload
6. SSE realtime connection
7. Logout flow

```

```
