import { z } from "zod";

export const loginSchema = z.object({
  email: z.string().email("Email không hợp lệ"),
  password: z.string().min(6, "Ít nhất 6 ký tự"),
});

export const updatePasswordSchema = z.object({
  email: z.string().email("Email không hợp lệ"),
  passwordOld: z.string().min(6, "Ít nhất 6 ký tự"),
  passwordNew: z.string().min(6, "Ít nhất 6 ký tự"),
  passwordConfirm: z.string().min(6, "Ít nhất 6 ký tự"),
});

export type UpdatePasswordForm = z.infer<typeof updatePasswordSchema>;

export type LoginForm = z.infer<typeof loginSchema>;
