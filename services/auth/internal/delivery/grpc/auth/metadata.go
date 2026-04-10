package authGrpc

import (
	"context"
	defErrors "libs/errors"
	"strings"

	"google.golang.org/grpc/metadata"
)

var AuthorizationHeader = "authorization"

func (h *Handler) getAuthToken(ctx context.Context) (string, error) {
	log := h.log.Named("getAuthToken")

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		log.Error(defErrors.MissingMetadata.Error())
		return "", defErrors.MissingMetadata
	}

	values := md.Get(AuthorizationHeader)
	if len(values) == 0 {
		log.Error(defErrors.MissingMetadataValue.Error())
		return "", defErrors.MissingMetadataValue
	}

	token := strings.TrimSpace(strings.TrimPrefix(values[0], "Bearer "))
	if token == "" {
		log.Error(MissingAuthToken.Error())
		return "", MissingAuthToken
	}

	return token, nil
}
