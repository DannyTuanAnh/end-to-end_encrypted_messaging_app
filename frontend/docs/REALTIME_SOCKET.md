## Realtime Socket (Socket.io)

This project uses `socket.io-client` for realtime messaging. Socket connection logic lives in `src/socket/socket.ts` and can be wrapped by a hook (`src/hooks/useSocket.ts`) that integrates with Zustand stores.

Connection

- Initialize the socket once (app entry) or lazily on login.
- Use the `VITE_SOCKET_URL` env variable for server address.

Example initialization (in `src/socket/socket.ts`)

```ts
import { io } from "socket.io-client";

export const socket = io(import.meta.env.VITE_SOCKET_URL, {
  autoConnect: false,
});

// Connect after login from a hook or auth flow:
// socket.auth = { token }; socket.connect();
```

Common events

- `send_message` — client -> server: send a new message payload.
- `receive_message` — server -> client: new message for a room or user.
- `typing` — client -> server: notify typing status.
- `user_online` / `user_offline` — server broadcasts presence updates.

Integration with state

- On `receive_message`, append to `chatStore` (e.g., `import { useChatStore } from 'src/store/chatStore'` and call the store updater).
- On `typing`, set a typing flag in `chatStore` keyed by conversation id.
- On `user_online`/`user_offline`, update `presence` state in `uiStore`/`chatStore`.

Example hook sketch

```ts
// Example hook (`src/hooks/useSocket.ts`)
function useSocket() {
  useEffect(() => {
    socket.on("receive_message", (msg) => {
      // update chat store
    });
    socket.on("typing", ({ userId, roomId }) => {
      // update typing state in store
    });
    return () => {
      socket.off("receive_message");
      socket.off("typing");
    };
  }, []);
}
```

Best practices

- Keep socket listeners idempotent and remove them when the component unmounts.
- Use namespaces or rooms on the server to limit event traffic.
- Guard large payloads and manage history via REST (fetch last N messages) rather than sending long arrays over sockets.
