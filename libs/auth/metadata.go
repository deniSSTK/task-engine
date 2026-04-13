package grpcAuth

import (
	"context"
	"strings"

	defErrors "github.com/deniSSTK/task-engine/libs/errors"
	"github.com/deniSSTK/task-engine/libs/logger"
	"google.golang.org/grpc/metadata"
)

const AuthorizationHeader = "authorization"

func ExtractAuthToken(ctx context.Context, log *logger.Logger) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
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
