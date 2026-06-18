package publishier

import (
	"context"

	"github.com/helium-mmm/Publishier/internal/domain"
)

type Publishier interface {
	Publish(
		ctx context.Context,
		post domain.Post,
		account domain.SocialAccount,
		botToken string,
	) (externalID string, err error)
}
