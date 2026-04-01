import { Link } from "react-router-dom";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

import { dataRoom } from "./dataRoom";

const userId = "u1"; // hardcoded current user ID for demo
const rooms = dataRoom.filter((r) => r.participants.includes(userId)); // in real app, fetch rooms for the user from server

export default function ChatList() {
  return (
    <aside className="h-full overflow-auto bg-card p-4">
      <div className="flex flex-col gap-3">
        {rooms.map((r) => (
          <Link to={`/chat/${r.id}`} key={r.id} className="no-underline">
            <div className="flex items-center gap-3 rounded-lg p-3 hover:bg-muted">
              <Avatar>
                <AvatarImage src={r.avatar} alt={r.name} />
                <AvatarFallback>{r.name.charAt(0)}</AvatarFallback>
              </Avatar>
              <div className="flex-1 min-w-0">
                <div className="font-medium truncate">{r.name}</div>
                <div className="flex justify-between items-center">
                  <div className="text-sm text-muted-foreground truncate">
                    {r.lastMessage}
                  </div>
                  {r.lastMessageAt && (
                    <div className="ml-2 text-xs text-muted-foreground">
                      {new Date(r.lastMessageAt).toLocaleTimeString()}
                    </div>
                  )}
                </div>
              </div>
            </div>
          </Link>
        ))}
      </div>
    </aside>
  );
}
