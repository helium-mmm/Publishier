import type { AuthResponse, Post, TelegramStatus } from "../types";

const BASE = "/api";

const EMAIL_RE = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

export function isValidEmail(email: string): boolean {
  return EMAIL_RE.test(email.trim());
}

class ApiError extends Error {
  constructor(
    message: string,
    readonly status: number,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

async function request<T>(
  path: string,
  options: RequestInit = {},
  token?: string | null,
): Promise<T> {
  const headers = new Headers(options.headers);
  headers.set("Content-Type", "application/json");
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const res = await fetch(`${BASE}${path}`, { ...options, headers });

  if (!res.ok) {
    let message = res.statusText;
    try {
      const body = (await res.json()) as { error?: string };
      if (body.error) message = body.error;
    } catch {
      // ignore
    }
    throw new ApiError(message, res.status);
  }

  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

export const api = {
  register: (email: string, password: string) =>
    request<AuthResponse>("/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    }),

  login: (email: string, password: string) =>
    request<AuthResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    }),

  getTelegramStatus: (token: string) =>
    request<TelegramStatus>("/accounts/telegram", { method: "GET" }, token),

  connectTelegram: (token: string, botToken: string, chatId: string) =>
    request<{ status: string }>(
      "/accounts/telegram",
      {
        method: "POST",
        body: JSON.stringify({ bot_token: botToken, chat_id: chatId }),
      },
      token,
    ),

  createPost: (token: string, content: string) =>
    request<Post>(
      "/posts",
      { method: "POST", body: JSON.stringify({ content }) },
      token,
    ),

  getPost: (token: string, id: string) =>
    request<Post>(`/posts/${id}`, { method: "GET" }, token),

  publishPost: (token: string, id: string) =>
    request<Post>(`/posts/${id}/publish`, { method: "POST" }, token),
};

export { ApiError };
