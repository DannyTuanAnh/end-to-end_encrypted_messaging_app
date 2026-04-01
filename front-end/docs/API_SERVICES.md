## API Services

Centralize HTTP logic in a service layer so components remain declarative.

Axios configuration

-- Create a single Axios instance (e.g., `src/services/axios.ts`) with defaults:

```ts
import axios from "axios";

// src/services/axios.ts
export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: { "Content-Type": "application/json" },
});
```

Token handling (auth)

- Store token in `useAuthStore` (Zustand). Do not store long-lived credentials in localStorage unless necessary. If persisted, prefer secure httpOnly cookies managed by backend.

- Use an interceptor on the `api` instance (`src/services/axios.ts`) to attach `Authorization: Bearer <token>` header:

```ts
api.interceptors.request.use((config) => {
  // example: import { getState } from '../store/authStore' or use a getter
  const token =
    (typeof window !== "undefined" && localStorage.getItem("token")) ||
    undefined;
  if (token) {
    config.headers = { ...config.headers, Authorization: `Bearer ${token}` };
  }
  return config;
});
```

Error handling

- Use a response interceptor to normalize errors and handle common status codes (401 -> token refresh or logout).
- Throw typed errors from service functions so callers can display messages or trigger retries.

API calling pattern

- Keep endpoint functions small and explicit. Example:

```ts
export const authService = {
  login: (data: LoginDto) => api.post("/auth/login", data),
  register: (data: RegisterDto) => api.post("/auth/register", data),
};
```

- Components should call hooks that wrap service calls (React Query `useMutation`/`useQuery`) rather than calling `api` directly.

Retries and loading

- Use React Query to manage retries, caching and invalidation. Use `isLoading` / `isFetching` to render spinners and `onError` handlers for notifications.

Security notes

- Avoid localStorage for tokens if backend supports httpOnly cookies.
- Sanitize inputs to API calls and validate responses.
