## State Management

This project uses Zustand for local/global client state and React Query for server-state (caching/fetching).

Zustand stores

- Pattern: create small focused stores per domain (authStore, uiStore, chatStore).
- Each store exports typed hooks or the store instance. Prefer selector usage for performance.

Example store shape (auth)

```ts
type AuthState = {
  user: User | null;
  token?: string;
  setUser: (u: User | null) => void;
  setToken: (t?: string) => void;
};

const useAuthStore = create<AuthState>()((set) => ({
  user: null,
  token: undefined,
  setUser: (user) => set({ user }),
  setToken: (token) => set({ token }),
}));
```

Global state structure (recommended)

- `authStore` — user, token, isAuthenticated helpers
- `uiStore` — sidebar open, active theme, modal states
- `chatStore` — current conversation id, typing state, message drafts

When to use global state vs local state

- Global state (Zustand):
  - cross-cutting concerns used in many places (auth, sidebar open)
  - short-lived UI state that must survive route changes (draft message)
- Local state (component useState):
  - purely presentational state scoped to a single component
  - form inputs that do not need to persist across routes

Data flow in the app

- User actions dispatch updates to stores and call service methods.
- For server data: React Query handles fetching and caching; components call queries/mutations and update UI on success.
- Socket events update Zustand stores directly (e.g., append received message to chatStore).

Best practices

- Keep stores small and domain-focused.
- Use selectors when reading from a store to avoid unnecessary renders.
- Avoid storing large arrays or binary data in Zustand; prefer React Query or indexed storage for large datasets.
