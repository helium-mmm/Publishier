package publishier

import (
	"context"
	"net/http"

	"github.com/helium-mmm/Publishier/internal/domain"
)

type TelegramPublishier struct {
	client *http.Client
	url string
}

func NewTelegramPublishier(client *http.Client, url string) *TelegramPublishier {
	return &TelegramPublishier{
		client: client,
		url: url,
	}
}

func (p *TelegramPublishier) Publish(
	ctx context.Context,
	post domain.Post,
	account domain.TelegramAccount,
) error {

}