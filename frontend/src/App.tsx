import { useCallback, useEffect, useState } from "react";
import { AuthForm } from "./components/AuthForm";
import { TelegramSetup } from "./components/TelegramSetup";
import { PostEditor } from "./components/PostEditor";
import { api } from "./api/client";
import type { Post } from "./types";

function App() {
  const [token, setToken] = useState<string | null>(
    () => localStorage.getItem("token"),
  );
  const [posts, setPosts] = useState<Post[]>([]);
  const [channelConnected, setChannelConnected] = useState(false);

  const refreshChannel = useCallback(async () => {
    if (!token) return;
    try {
      const status = await api.getTelegramStatus(token);
      setChannelConnected(status.connected);
    } catch {
      setChannelConnected(false);
    }
  }, [token]);

  useEffect(() => {
    refreshChannel();
  }, [refreshChannel]);

  function handleLogout() {
    localStorage.removeItem("token");
    setToken(null);
    setPosts([]);
    setChannelConnected(false);
  }

  if (!token) {
    return (
      <div className="app">
        <div className="glow glow-1" />
        <div className="glow glow-2" />
        <div className="container">
          <header className="brand">
            <span className="logo">P</span>
            <h1>Publishier</h1>
          </header>
          <p className="tagline">Публикуйте посты в Telegram в один клик</p>
          <AuthForm onAuth={setToken} />
        </div>
      </div>
    );
  }

  return (
    <div className="app">
      <div className="glow glow-1" />
      <div className="glow glow-2" />
      <div className="container">
        <header className="topbar">
          <div className="brand-inline">
            <span className="logo sm">P</span>
            <h1>Publishier</h1>
          </div>
          <button type="button" className="ghost" onClick={handleLogout}>
            Выйти
          </button>
        </header>

        <section className="card">
          <TelegramSetup
            token={token}
            onStatusChange={setChannelConnected}
          />
        </section>

        <section className="card">
          <p className="section-title">Редактор</p>
          <PostEditor
            token={token}
            posts={posts}
            channelConnected={channelConnected}
            onPostCreated={(post) => setPosts((prev) => [post, ...prev])}
            onPostUpdated={(post) =>
              setPosts((prev) => prev.map((p) => (p.id === post.id ? post : p)))
            }
          />
        </section>
      </div>
    </div>
  );
}

export default App;
