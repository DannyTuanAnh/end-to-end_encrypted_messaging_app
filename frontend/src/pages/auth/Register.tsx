import React from "react";
import { useNavigate, Link } from "react-router-dom";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

export default function Register() {
  const navigate = useNavigate();

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    // fake register: set token
    localStorage.setItem("token", "demo-token");
    navigate("/");
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <h2 className="text-xl font-semibold">Create account</h2>
      <Input placeholder="Name" required />
      <Input placeholder="Email" type="email" required />
      <Input placeholder="Password" type="password" required />
      <div className="flex items-center justify-between">
        <Button type="submit">Register</Button>
        <Link to="/auth/login" className="text-sm text-primary">
          Sign in
        </Link>
      </div>
    </form>
  );
}
