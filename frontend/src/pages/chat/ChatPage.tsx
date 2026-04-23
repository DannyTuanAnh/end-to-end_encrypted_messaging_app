import { useState } from "react";
import { useParams } from "react-router-dom";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { useChat } from "@/hooks/chat/useChat";
import { useMessages } from "@/hooks/chat/useMessages";
import QRCode from "@/components/common/QRCode";

export default function ChatPage() {
  const { id } = useParams();
  const uid = "user_123456"; // hardcoded user ID for demo
  const currenUser = "u1"; // hardcoded current user ID for demo
  const { sendMessage } = useChat(currenUser);
  const messages = useMessages(String(id)); // reactive messages for the room
  const [text, setText] = useState("");
  if (id === undefined)
    return (
      <div className="h-full flex flex-col items-center justify-center gap-4">
        <p className="text-muted-foreground">
          Select a chat to start messaging
        </p>
        <h3 className="text-lg font-medium">Your QR Code</h3>
        <p className="text-sm text-muted-foreground">
          Scan this QR code to add me as a contact.
        </p>
        <QRCode uid={uid} />
      </div>
    );
  return (
    <div className="h-full flex-1 flex flex-col gap-4 overflow-hidden">
      <div>
        <h1 className="text-2xl font-semibold">Chat #{id}</h1>
      </div>

      <div className="w-full h-full flex-1 overflow-y-auto flex flex-col gap-2 p-4 border rounded-lg">
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
              className={`p-4 ${m.senderId === currenUser ? "bg-primary text-white w-1/2 ml-auto" : "bg-muted text-default w-1/2"}`}
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
          <Input
            placeholder="Type a message"
            value={text}
            onChange={(e) => setText(e.target.value)}
          />
          <Button
            onClick={() => {
              if (!id || !text.trim()) return;
              sendMessage(String(id), currenUser, text.trim());
              setText("");
            }}
          >
            Send
          </Button>
        </div>
      </div>
    </div>
  );
}
