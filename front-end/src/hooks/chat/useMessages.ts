import { useChatStore } from "@/stores/chatStore";

export function useMessages(roomId?: string) {
  const messages = useChatStore((s) => s.messages);
  return roomId ? messages[roomId] || [] : [];
}
