import { useMemo } from "react";
import { useChatStore } from "@/stores/chatStore";

export function useChat(currentUserId: string) {
  const rooms = useChatStore((s) => s.rooms);
  const sendMessage = useChatStore((s) => s.sendMessage);

  const myRooms = useMemo(() => {
    return rooms.filter((r) => r.participants.includes(currentUserId));
  }, [rooms, currentUserId]);

  return {
    rooms: myRooms,
    allRooms: rooms,
    sendMessage,
  };
}
