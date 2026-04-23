## Future Improvements

Suggested improvements

- Authentication: implement refresh token flow and migrate to httpOnly cookies for improved security.
- Message history: implement pagination and server-side cursors to avoid loading all messages at once.

Performance optimizations

- Use virtualization (react-window / react-virtual) for long message lists.
- Lazy-load route components and large UI bundles to reduce initial bundle size.
- Memoize heavy renderers and use selectors for Zustand to prevent re-renders.

Scalability suggestions

- Move socket event handling into a dedicated service that can be tested separately.
- Introduce a normalized client cache for messages (like RTK Query normalization or a small in-memory index) to manage message updates and deduplication.
- Add end-to-end tests for critical flows (login, send/receive message).

Observability

- Add error reporting (Sentry) and performance monitoring to capture runtime issues.

Developer experience

- Add `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md` and minimal setup scripts.
- Add storybook or a component playground for shared components.
