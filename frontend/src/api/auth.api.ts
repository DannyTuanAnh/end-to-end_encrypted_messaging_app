import { api } from "@/lib/axios";

export type LoginGoogleBody = {
  auth_code: string;
};

export type User = {
  uuid: string;
  email?: string;
  name?: string;
  avatar_url?: string;
};

export type ApiResponse<T> = {
  success: boolean;
  data: T;
  message?: string;
};

export async function loginGoogle(body: LoginGoogleBody) {
  const response = await api.post("/auth/google/login", body);

  return response.data;
}

export async function logoutApi() {
  const response = await api.delete("/auth/logout");

  return response.data;
}

export async function getProfile() {
  const response = await api.get<ApiResponse<User>>("/user/profile");

  return response.data.data;
}
