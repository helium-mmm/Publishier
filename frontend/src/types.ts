export interface Post {
  id: string;
  content: string;
  status: "DRAFT" | "PUBLISHED" | "FAILED";
  created_at: string;
  published_at?: string;
}

export interface AuthResponse {
  token: string;
}

export interface TelegramStatus {
  connected: boolean;
  chat_id?: string;
  platform?: string;
}
