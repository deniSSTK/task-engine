package grpcAuth

import (
	"context"

	commonv1 "github.com/deniSSTK/task-engine/gen/proto/common/v1"
	defErrors "github.com/deniSSTK/task-engine/libs/errors"
	userDomain "github.com/deniSSTK/task-engine/libs/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type authUserCtxKey struct{}

func WithAuthUser(ctx context.Context, authUser *userDomain.AuthUser) context.Context {
	return context.WithValue(ctx, authUserCtxKey{}, authUser)
}

func AuthUserFromContext(ctx context.Context) (*userDomain.AuthUser, bool) {
	authUser, ok := ctx.Value(authUserCtxKey{}).(*userDomain.AuthUser)
	return authUser, ok
}

func UnaryAuthInterceptor(
	registry *PolicyRegistry,
	verifier AuthVerifier,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		log := registry.log.Named("UnaryAuthInterceptor")

		policy := registry.MustGet(info.FullMethod) // TODO: I think it's not a good idea to panic here

		if policy == commonv1.AuthPolicy_AUTH_POLICY_PUBLIC {
			return handler(ctx, req)
		}

		if _, ok := metadata.FromIncomingContext(ctx); !ok {
			log.Error("missing metadata")
			return nil, registry.errorWrapper.Unauthenticated(defErrors.MissingMetadata)
		}

		authUser, err := verifier.Verify(ctx)
		if err != nil {
			log.Error("failed to verify", zap.Error(err))
			return nil, registry.errorWrapper.Unauthenticated(err)
		}

		if policy == commonv1.AuthPolicy_AUTH_POLICY_ADMIN &&
			authUser.Role != userDomain.RoleUser {
			log.Error("admin only")
			return nil, registry.errorWrapper.PermissionDenied()
		}

		ctx = WithAuthUser(ctx, authUser)

		return handler(ctx, req)
	}
}
