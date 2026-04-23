import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { loginSchema } from "@/lib/schema";
import type { LoginForm } from "@/lib/schema";
import { useNavigate, Link } from "react-router-dom";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useAuthContext } from "@/context/AuthContext";

export default function Login() {
  const navigate = useNavigate();
  const { login } = useAuthContext();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
    mode: "onChange",
  });

  function onSubmit(data: LoginForm) {
    // demo login via auth hook
    login(data.email, data.password);
    navigate("/");
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
      <h2 className="text-xl font-semibold">Sign in</h2>
      <Input {...register("email")} placeholder="Email" type="email" required />
      {errors.email && (
        <p className="text-red-500 text-sm">{errors.email.message}</p>
      )}
      <Input
        {...register("password")}
        placeholder="Password"
        type="password"
        required
      />
      {errors.password && (
        <p className="text-red-500 text-sm">{errors.password.message}</p>
      )}
      <div className="flex items-center justify-between">
        <Button type="submit" disabled={isSubmitting}>
          Sign in
        </Button>
        <Link to="/auth/register" className="text-sm text-primary">
          Create account
        </Link>
      </div>
    </form>
  );
}
