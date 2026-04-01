import { useParams } from "react-router-dom";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { dataMessage } from "./dataMessage";

export default function ChatPage() {
  const { id } = useParams();
  const currenUser = "u1"; // hardcoded current user ID for demo
  const messages = dataMessage[String(id)]; // in real app, fetch messages for the room from server
  if (id === undefined)
    return (
      <div className="h-full flex items-center justify-center">
        <p className="text-muted-foreground">
          Select a chat to start messaging
        </p>
      </div>
    );
  return (
    <div className="h-full flex flex-col gap-4">
      <div>
        <h1 className="text-2xl font-semibold">Chat #{id}</h1>
      </div>

      <div className="w-full flex-1 overflow-auto flex flex-col gap-3 p-4 border rounded-lg">
        {messages?.map((m) => (
          <div key={m.id} className={`p-3 rounded-lg flex space-x-2`}>
            {m.senderId !== currenUser && (
              <Avatar>
                <AvatarFallback>
                  {m.senderId.charAt(0).toUpperCase()}
                </AvatarFallback>
              </Avatar>
            )}

            <Card
              className={`p-4 mb-4 ${m.senderId === currenUser ? "bg-primary text-white w-1/2 ml-auto" : "bg-muted text-default w-1/2"}`}
            >
              <p className="text-base">{m.content}</p>
              <p className="text-xs">
                {m.createdAt && new Date(m.createdAt).toLocaleTimeString()}
              </p>
            </Card>
          </div>
        ))}
      </div>

      <div className="mt-2">
        <div className="flex gap-2">
          <Input placeholder="Type a message" />
          <Button>Send</Button>
        </div>
      </div>
    </div>
  );
}
