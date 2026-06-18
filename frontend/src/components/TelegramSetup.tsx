import { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import { api, ApiError } from "../api/client";
import type { TelegramStatus } from "../types";

interface Props {
  token: string;
  onStatusChange?: (connected: boolean) => void;
}

export function TelegramSetup({ token, onStatusChange }: Props) {
  const [status, setStatus] = useState<TelegramStatus | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [botToken, setBotToken] = useState("");
  const [chatId, setChatId] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function loadStatus() {
    try {
      const data = await api.getTelegramStatus(token);
      setStatus(data);
      onStatusChange?.(data.connected);
    } catch {
      setStatus({ connected: false });
      onStatusChange?.(false);
    }
  }

  useEffect(() => {
    loadStatus();
  }, [token]);

  useEffect(() => {
    if (!modalOpen) return;
    const prev = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => {
      document.body.style.overflow = prev;
    };
  }, [modalOpen]);

  function openModal() {
    setBotToken("");
    setChatId(status?.chat_id ?? "");
    setError("");
    setModalOpen(true);
  }

  function closeModal() {
    setModalOpen(false);
    setBotToken("");
    setError("");
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await api.connectTelegram(token, botToken, chatId);
      await loadStatus();
      closeModal();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Ошибка");
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      <div className="channel-card">
        <div className="channel-info">
          <span className={`channel-dot ${status?.connected ? "online" : ""}`} />
          <div>
            <p className="channel-label">Telegram-канал</p>
            {status?.connected ? (
              <p className="channel-value">{status.chat_id}</p>
            ) : (
              <p className="channel-value muted">Не подключён</p>
            )}
          </div>
        </div>
        <button type="button" className="ghost" onClick={openModal}>
          {status?.connected ? "Изменить" : "Подключить"}
        </button>
      </div>

      {modalOpen &&
        createPortal(
          <div className="modal-overlay" onClick={closeModal}>
            <div className="modal" onClick={(e) => e.stopPropagation()} role="dialog" aria-modal="true">
              <div className="modal-header">
                <h2>{status?.connected ? "Изменить канал" : "Подключить Telegram"}</h2>
                <button type="button" className="modal-close" onClick={closeModal}>
                  ×
                </button>
              </div>
              <p className="modal-hint">
                Укажите bot token и chat id канала, куда будут уходить посты.
              </p>
              <form onSubmit={handleSubmit}>
                <div className="field">
                  <label htmlFor="bot-token">Bot Token</label>
                  <input
                    id="bot-token"
                    type="password"
                    value={botToken}
                    onChange={(e) => setBotToken(e.target.value)}
                    placeholder="123456:ABC..."
                    required
                    autoFocus
                  />
                </div>
                <div className="field">
                  <label htmlFor="chat-id">Chat ID</label>
                  <input
                    id="chat-id"
                    type="text"
                    value={chatId}
                    onChange={(e) => setChatId(e.target.value)}
                    placeholder="@channel или -100..."
                    required
                  />
                </div>
                {error && <p className="error">{error}</p>}
                <div className="modal-actions">
                  <button type="button" onClick={closeModal}>
                    Отмена
                  </button>
                  <button type="submit" className="primary" disabled={loading}>
                    {loading ? "Сохранение..." : "Сохранить"}
                  </button>
                </div>
              </form>
            </div>
          </div>,
          document.body,
        )}
    </>
  );
}
