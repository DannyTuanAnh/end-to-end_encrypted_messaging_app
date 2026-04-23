## Coding Conventions

Naming conventions

- Components: PascalCase (e.g., `ChatList`, `MessageItem`).
- Hooks: `use` prefix, camelCase (e.g., `useSocket`, `useMobile`).
- Stores: `useXStore` (e.g., `useAuthStore`).
- Services: camelCase with domain (e.g., `authService`, `chatService`).

File naming

- Use `.tsx` for files that return JSX, `.ts` for plain logic.
- Group component files with same name (`MyComponent.tsx`) and optional `styles.css` or `index.ts` barrel file.

Component structure

- Prefer small, focused components. Each component should:
  - Accept props and be as stateless as possible.
  - Keep side-effects in hooks or higher-order components.

Hook rules

- Keep hooks pure and prefix with `use`.
- Return stable references (use `useCallback` / `useMemo`) where consumers depend on function identity.

Service pattern

- Centralize Axios instance and export small functions for endpoints.
- Avoid network calls inside presentational components — use hooks (React Query) to call services.

Store pattern

- Keep stores focused and export selector helpers.
- Avoid mutating state directly — use setter functions.

TypeScript usage

- Use explicit types for public interfaces (props, store state, service responses).
- Keep `any` to a minimum; prefer `unknown` and narrow types.

Formatting and linting

- Run Prettier and ESLint configured in the repo. Keep files formatted and lint warnings addressed before PRs.
