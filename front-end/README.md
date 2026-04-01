# Chat App — Frontend

Project overview

- A modern chat frontend built with React + TypeScript, Vite, TailwindCSS and shadcn/ui. It connects to a REST API via Axios and to realtime messaging via Socket.io. State is managed with Zustand and server state with React Query.

Tech stack

- React 18 + TypeScript
- Vite
- TailwindCSS
- shadcn/ui components
- React Router
- Zustand (global state)
- React Query (server state)
- Axios (HTTP client)
- Socket.io-client (realtime)

Installation

1. Clone the repo and change into the frontend folder:

```bash
git clone <repo>
cd front-end
```

2. Install dependencies:

```bash
npm install
# or
pnpm install
```

Run project (development)

```bash
npm run dev
# opens Vite dev server (hot reload)
```

Build project

```bash
npm run build
npm run preview
```

Environment variables

- Create a `.env` file in `front-end/` or use environment-specific files.
- Typical variables the app expects:

- `VITE_API_BASE_URL` — base URL for the REST API
- `VITE_SOCKET_URL` — URL for socket.io server
- `VITE_APP_NAME` — optional display name

Folder structure overview (top-level)

- `src/` — application source
- `src/components/` — shared UI components (shadcn wrappers + custom controls)
- `src/pages/` — route-level pages (Home, Profile, auth)
- `src/hooks/` — custom React hooks
- `src/lib/` — small helpers and utilities
- `src/store/` — Zustand stores
- `src/style/` — CSS / Tailwind entry (global styles)
- `src/router.ts` — app routing configuration
- `src/main.tsx` — app entry

For more details see `FOLDER_STRUCTURE.md` and per-topic docs in this folder.
