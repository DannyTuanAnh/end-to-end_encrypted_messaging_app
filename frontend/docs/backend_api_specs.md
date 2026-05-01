# Backend API Specifications

Dưới đây là các thông tin được tổng hợp từ mã nguồn backend (API Gateway).

## 1. Request Body Format

Các API nhận dữ liệu dưới dạng JSON (`application/json`) hoặc Multipart Form Data (`multipart/form-data`) tùy thuộc vào loại dữ liệu.

Ví dụ:

- **LoginGoogle**:
  ```json
  {
    "auth_code": "string"
  }
  ```
- **UpdateProfile** (`PUT /api/v1/user/profile` - `multipart/form-data`):
  - `avatar_url`: File (ảnh)
  - `name`: string (optional)
  - `birthday`: string (optional)
  - `phone`: string (optional)
- **VerifyIDTokenOTP**:
  ```json
  {
    "id_token": "string"
  }
  ```

## 2. Response Mẫu

- **Thành công (Success):**

  ```json
  {
    "success": true,
    "data": { ... } // (Optional) chứa dữ liệu payload
  }
  ```

  Hoặc có thêm `message`:

  ```json
  {
    "success": true,
    "message": "Thông báo thành công"
  }
  ```

- **Thất bại (Error):**
  ```json
  {
    "error": "Nội dung lỗi",
    "code": "ERROR_CODE",
    "details": "Chi tiết lỗi (Optional)"
  }
  ```

## 3. Auth Flow

1. **Đăng nhập**: Client gửi yêu cầu đến `/api/v1/auth/google/login` với body chứa `auth_code`.
2. **Cấp phiên (Session)**: Backend gọi Auth Service, sau đó thiết lập hai HTTP-Only & Secure Cookies là `session_id` và `device_id`.
3. **Lưu trữ**: Dữ liệu phiên đăng nhập được lưu trữ trong Redis và Database.
4. **Xác thực**: Middleware `AuthMiddleware` sẽ lấy `session_id` và `device_id` từ cookies ở các request tiếp theo, sau đó đối chiếu thông tin với Redis (hoặc Database nếu Redis không có). Nếu hợp lệ, `user_id` và `user_uuid` sẽ được gán vào context.

_(Ghi chú: Đối với các API công khai (`/api/v1/`), hệ thống còn sử dụng `ApiKeyMiddleware` yêu cầu header `X-Api-Key`)_.

## 4. Token Format

Hệ thống không sử dụng JWT Bearer Token trên frontend. Thay vào đó, **token được lưu dưới dạng UUID string qua HTTP-only cookies**:

- `session_id`: Chuỗi UUID (vd: `123e4567-e89b-12d3-a456-426614174000`)
- `device_id`: Chuỗi UUID

## 5. SSE Auth (Server-Sent Events)

- **Endpoint**: `GET /api/v1/notify/sse`
- **Xác thực**: Flow xác thực cho SSE sử dụng chung cơ chế **Cookie `session_id` và `device_id`** qua `AuthMiddleware` giống hệt các REST API khác. Nếu Cookie hợp lệ, SSE connection sẽ được thiết lập thành công.

## 6. Upload Avatar API

Upload Avatar được xử lý tích hợp vào API **UpdateProfile**:

- **Route**: `PUT /api/v1/user/profile` (Nằm trong route group bảo vệ).
- **Body**: Sử dụng `multipart/form-data`. Gửi file ảnh qua key `avatar_url`.
- **Cơ chế**: File upload được tải trực tiếp lên Google Cloud Storage (Bucket `GOOGLE_CLOUD_STORAGE_BUCKET_RAW`).
- Có một API khác là `POST /api/v1/user/report-avatar` (body: `uuid` và `name`) để ghi nhận sau khi xử lý avatar.

## 7. Message APIs

Hiện tại ở lớp API Gateway (Backend này), **chưa có endpoints nào liên quan đến việc gửi/nhận tin nhắn (Message APIs) được thiết lập** (Các module hiện tại chỉ bao gồm `AuthModule`, `UserModule`, và `NotifyModule`).

## 8. Websocket/SSE Event Format

Mỗi sự kiện trả về qua Server-Sent Events (SSE) tuân theo chuẩn SSE:

- **Kết nối thành công**:
  ```text
  : connected\n\n
  ```
- **Ping/Heartbeat** (gửi mỗi 5s để giữ kết nối):
  ```text
  : ping\n\n
  ```
- **Sự kiện dữ liệu từ Redis Pub/Sub**:
  ```text
  data: {"user_id":"...","status":"...","file_path":"..."}\n\n
  ```
  Trong đó payload là một chuỗi JSON gồm:
  ```json
  {
    "user_id": "string",
    "status": "string",
    "file_path": "string"
  }
  ```

## 9. Status Code Convention

Quy ước custom error code (bên trong JSON response) ánh xạ sang HTTP Status Code:

- `BAD_REQUEST` -> `400 Bad Request`
- `UNAUTHORIZED` -> `401 Unauthorized`
- `NOT_FOUND` -> `404 Not Found`
- `CONFLICT` -> `409 Conflict`
- `TOO_MANY_REQUESTS` -> `429 Too Many Requests`
- `INTERNAL_SERVER_ERROR` -> `500 Internal Server Error`

Khi xử lý lỗi từ gRPC service, mã lỗi của gRPC (như `InvalidArgument`, `PermissionDenied`, `Unauthenticated`) cũng sẽ được ánh xạ sang các HTTP Status Code tương ứng (400, 403, 401...).

## 10. Validation Error Format

Định dạng lỗi khi validation đầu vào thất bại có 2 dạng chính:

- **Lỗi validation do custom gin validator (HTTP Layer)**:
  ```json
  {
    "error": {
      "field_name": "Câu thông báo lỗi (vd: field_name là bắt buộc)",
      "another_field": "Câu thông báo lỗi"
    }
  }
  ```
- **Lỗi validation do gRPC (Protobuf Validation)**:
  ```json
  {
    "error": "VALIDATION_ERROR",
    "detail": {
      "field_name": "Error description"
    }
  }
  ```
