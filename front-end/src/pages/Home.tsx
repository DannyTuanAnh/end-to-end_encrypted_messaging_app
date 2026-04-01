import { Link } from "react-router-dom";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function Home() {
  return (
    <div className="max-w-3xl min-w-sm w-full">
      <h1 className="text-2xl font-semibold mb-4">Welcome</h1>

      <Card className="p-4 mb-4">
        <p className="mb-3">This is the home page. Quick links:</p>
        <div className="flex gap-2">
          <Link to="/chat">
            <Button>Open Chat</Button>
          </Link>
          <Link to="/profile">
            <Button variant="outline">Profile</Button>
          </Link>
        </div>
      </Card>
    </div>
  );
}
