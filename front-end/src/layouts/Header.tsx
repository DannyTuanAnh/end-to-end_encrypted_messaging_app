import { Link, useNavigate } from "react-router-dom";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";

export default function Header() {
  const navigate = useNavigate();

  function handleLogout() {
    localStorage.removeItem("token");
    navigate("/auth/login");
  }

  return (
    <header className="h-16 flex items-center justify-between gap-4 border-b px-4 py-3">
      <div className="flex items-center gap-3">
        <Link to="/" className="text-lg font-semibold">
          Chat App
        </Link>
      </div>

      <div className="flex items-center gap-3">
        <Link to="/profile" aria-label="Profile">
          <Avatar>
            <AvatarImage src="/assets/avatar-placeholder.png" alt="User" />
            <AvatarFallback>U</AvatarFallback>
          </Avatar>
        </Link>
        <Button variant="ghost" onClick={handleLogout}>
          Sign out
        </Button>
      </div>
    </header>
  );
}
