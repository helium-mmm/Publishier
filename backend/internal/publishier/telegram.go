package publishier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/helium-mmm/Publishier/internal/domain"
)

type TelegramPublishier struct {
	client  *http.Client
	baseURL string
}

func NewTelegramPublishier(client *http.Client, baseURL string) *TelegramPublishier {
	return &TelegramPublishier{
		client:  client,
		baseURL: baseURL,
	}
}

type sendMessageRequest struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type telegramResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
	} `json:"result"`
	Description string `json:"description"`
}

func (p *TelegramPublishier) Publish(
	ctx context.Context,
	post domain.Post,
	account domain.SocialAccount,
	botToken string,
) (string, error) {
	body, err := json.Marshal(sendMessageRequest{
		ChatID: account.ChatID,
		Text:   post.Content,
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", p.baseURL, botToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tgResp telegramResponse
	if err := json.Unmarshal(respBody, &tgResp); err != nil {
		return "", err
	}

	if !tgResp.OK {
		return "", fmt.Errorf("telegram api error: %s", tgResp.Description)
	}

	return strconv.Itoa(tgResp.Result.MessageID), nil
}
