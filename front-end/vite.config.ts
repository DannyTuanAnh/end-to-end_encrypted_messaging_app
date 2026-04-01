// import { defineConfig } from 'vite'
// import react from '@vitejs/plugin-react'

// // https://vite.dev/config/
// export default defineConfig({
//   plugins: [react()],
// })

import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

// import fs from "fs";
// import path from "path";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  // server: {
  //   https: {
  //     key: fs.readFileSync(
  //       path.resolve(
  //         "../../end-to-end_encrypted_messaging_app/internal/certs/localhost+2-key.pem",
  //       ), // đường dẫn key
  //     ),
  //     cert: fs.readFileSync(
  //       path.resolve(
  //         "../../end-to-end_encrypted_messaging_app/internal/certs/localhost+2.pem",
  //       ), // đường dẫn cert
  //     ),
  //   },
  //   port: 5173,
  // },
});
