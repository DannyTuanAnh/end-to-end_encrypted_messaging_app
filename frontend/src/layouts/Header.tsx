import { Link, useNavigate } from "react-router-dom";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { MessageCircle, LogOut } from "lucide-react";

export default function Header() {
  const navigate = useNavigate();

  function handleLogout() {
    localStorage.removeItem("token");
    navigate("/auth/login");
  }

  return (
    <header className="h-16 bg-background fixed top-0 left-0 right-0 flex items-center justify-between gap-4 border-b px-4 py-3 z-10 over-flow-hidden">
      <div className="flex items-center gap-3">
        <Link to="/" className="text-lg font-semibold">
          Chat App
        </Link>
      </div>

      <div className="flex items-center gap-3">
        <Link to="/chat" className="text-sm">
          <Button variant="default">
            <MessageCircle className="h-4 w-4" />
            Chat
          </Button>
        </Link>
        <Link to="/profile" aria-label="Profile">
          <Avatar>
            <AvatarImage src="/assets/avatar-placeholder.png" alt="User" />
            <AvatarFallback>U</AvatarFallback>
          </Avatar>
        </Link>
        <Button variant="destructive" onClick={handleLogout}>
          <LogOut className="h-4 w-4" />
          Logout
        </Button>
      </div>
    </header>
  );
}
