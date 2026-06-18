import { useState } from "react";
import { api, ApiError, isValidEmail } from "../api/client";

interface Props {
  onAuth: (token: string) => void;
}

export function AuthForm({ onAuth }: Props) {
  const [mode, setMode] = useState<"login" | "register">("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    if (!isValidEmail(email)) {
      setError("Введите корректный email");
      return;
    }

    setLoading(true);
    try {
      const fn = mode === "login" ? api.login : api.register;
      const { token } = await fn(email.trim().toLowerCase(), password);
      localStorage.setItem("token", token);
      onAuth(token);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Ошибка");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="card">
      <div className="tabs">
        <button
          type="button"
          className={mode === "login" ? "active" : ""}
          onClick={() => setMode("login")}
        >
          Вход
        </button>
        <button
          type="button"
          className={mode === "register" ? "active" : ""}
          onClick={() => setMode("register")}
        >
          Регистрация
        </button>
      </div>

      <form onSubmit={handleSubmit}>
        <div className="field">
          <label htmlFor="email">Email</label>
          <input
            id="email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="you@example.com"
            required
            autoComplete="email"
          />
        </div>
        <div className="field">
          <label htmlFor="password">Пароль</label>
          <input
            id="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            autoComplete={mode === "login" ? "current-password" : "new-password"}
          />
        </div>
        <button type="submit" className="primary" disabled={loading}>
          {loading ? "Загрузка..." : mode === "login" ? "Войти" : "Создать аккаунт"}
        </button>
        {error && <p className="error">{error}</p>}
      </form>
    </div>
  );
}
