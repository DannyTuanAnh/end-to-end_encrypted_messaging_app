import { useMemo } from "react";
import { useChatStore } from "@/stores/chatStore";

export function useChat(currentUserId: string | null) {
  if (!currentUserId) {
    throw new Error("currentUserId is required for useChat");
  }
  const rooms = useChatStore((s) => s.rooms);
  const sendMessage = useChatStore((s) => s.sendMessage);

  const myRooms = useMemo(() => {
    return rooms.filter((room) => room.participants.includes(currentUserId));
  }, [rooms, currentUserId]);

  return {
    rooms: myRooms,
    allRooms: rooms,
    sendMessage,
  };
}
