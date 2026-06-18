import { useState } from "react";
import type { Post } from "../types";
import { api, ApiError } from "../api/client";

interface Props {
  token: string;
  posts: Post[];
  channelConnected: boolean;
  onPostCreated: (post: Post) => void;
  onPostUpdated: (post: Post) => void;
}

export function PostEditor({
  token,
  posts,
  channelConnected,
  onPostCreated,
  onPostUpdated,
}: Props) {
  const [content, setContent] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [publishingId, setPublishingId] = useState<string | null>(null);

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!content.trim()) return;
    setError("");
    setLoading(true);
    try {
      const post = await api.createPost(token, content.trim());
      onPostCreated(post);
      setContent("");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Ошибка");
    } finally {
      setLoading(false);
    }
  }

  async function handlePublish(id: string) {
    if (!channelConnected) {
      setError("Сначала подключите Telegram-канал в настройках");
      return;
    }

    setError("");
    setPublishingId(id);
    try {
      const post = await api.publishPost(token, id);
      onPostUpdated(post);
    } catch (err) {
      if (err instanceof ApiError) {
        if (err.status === 400 && err.message.includes("connected")) {
          setError("Сначала подключите Telegram-канал в настройках");
        } else {
          setError(err.message);
        }
      } else {
        setError("Ошибка публикации");
      }
    } finally {
      setPublishingId(null);
    }
  }

  return (
    <div>
      <form onSubmit={handleCreate}>
        <div className="field">
          <label htmlFor="content">Текст поста</label>
          <textarea
            id="content"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="Напишите пост..."
            required
          />
        </div>
        <button type="submit" className="primary" disabled={loading}>
          {loading ? "Сохранение..." : "Сохранить черновик"}
        </button>
      </form>

      {error && <p className="error">{error}</p>}

      <div className="posts-section">
        <p className="section-title">Посты</p>
        {posts.length === 0 ? (
          <p className="empty">Пока нет постов — напишите первый</p>
        ) : (
          posts.map((post) => (
            <div key={post.id} className="post-item">
              <p>{post.content}</p>
              <div className="post-meta">
                <span className={`badge ${post.status.toLowerCase()}`}>
                  {post.status}
                </span>
                {post.status === "DRAFT" && (
                  <button
                    type="button"
                    className="primary sm"
                    disabled={publishingId === post.id}
                    onClick={() => handlePublish(post.id)}
                  >
                    {publishingId === post.id ? "..." : "Опубликовать"}
                  </button>
                )}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
