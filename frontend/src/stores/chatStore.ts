import { create } from "zustand";
import type { Room } from "@/data/dataRoom";
import type { Message } from "@/data/dataMessage";
import { dataRoom } from "@/data/dataRoom";
import { dataMessage } from "@/data/dataMessage";

type MessagesMap = Record<string, Message[]>;

type ChatState = {
  rooms: Room[];
  messages: MessagesMap;

  setRooms: (rooms: Room[]) => void;
  sendMessage: (roomId: string, senderId: string, content: string) => void;
};

export const useChatStore = create<ChatState>((set) => ({
  rooms: dataRoom,
  messages: dataMessage,

  setRooms: (rooms) => set({ rooms }),

  sendMessage: (roomId, senderId, content) => {
    const newMessage: Message = {
      id: `m${Date.now()}`,
      roomId,
      senderId,
      content,
      createdAt: new Date().toISOString(),
    };

    set((state) => ({
      messages: {
        ...state.messages,
        [roomId]: [...(state.messages[roomId] ?? []), newMessage],
      },
      rooms: state.rooms.map((room) =>
        room.id === roomId
          ? {
              ...room,
              lastMessage: content,
              lastMessageAt: newMessage.createdAt,
            }
          : room,
      ),
    }));
  },
}));
