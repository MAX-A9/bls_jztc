/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_APP_NAME: string;
  readonly VITE_API_BASE_URL: string;
  readonly VITE_APP_ENV: string;
  readonly VITE_API_SERVER: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
} 