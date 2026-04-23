// Global Type Declarations for window.google
export {};

declare global {
  interface Window {
    // simplified shape for Google Identity Services used in this app
    google?: {
      accounts?: {
        oauth2?: {
          initCodeClient?: (opts?: any) => { requestCode?: () => void } & any;
        };
      };
    };
  }
}
