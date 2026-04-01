import React from "react";
import { useNavigate, Link } from "react-router-dom";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

export default function Login() {
  const navigate = useNavigate();

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    // fake login: set token
    localStorage.setItem("token", "demo-token");
    navigate("/");
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <h2 className="text-xl font-semibold">Sign in</h2>
      <Input placeholder="Email" type="email" required />
      <Input placeholder="Password" type="password" required />
      <div className="flex items-center justify-between">
        <Button type="submit">Sign in</Button>
        <Link to="/auth/register" className="text-sm text-primary">
          Create account
        </Link>
      </div>
    </form>
  );
}
