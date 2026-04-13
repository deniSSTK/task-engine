package grpcAuth

import (
	"context"

	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
	commonv1 "github.com/deniSSTK/task-engine/gen/proto/common/v1"
	defErrors "github.com/deniSSTK/task-engine/libs/errors"
	userDomain "github.com/deniSSTK/task-engine/libs/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type authUserCtxKey struct{}

func WithAuthUser(ctx context.Context, authUser *authv1.AuthUser) context.Context {
	return context.WithValue(ctx, authUserCtxKey{}, authUser)
}

func AuthUserFromContext(ctx context.Context) (*authv1.AuthUser, bool) {
	authUser, ok := ctx.Value(authUserCtxKey{}).(authv1.AuthUser)
	return &authUser, ok
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
		policy := registry.MustGet(info.FullMethod)

		if policy == commonv1.AuthPolicy_AUTH_POLICY_PUBLIC {
			return handler(ctx, req)
		}

		if _, ok := metadata.FromIncomingContext(ctx); !ok {
			return nil, status.Error(codes.Unauthenticated, defErrors.MissingMetadata.Error())
		}

		authUser, err := verifier.Verify(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		if policy == commonv1.AuthPolicy_AUTH_POLICY_ADMIN &&
			userDomain.UserRole(authUser.Role) != userDomain.Admin {
			return nil, status.Error(codes.Unauthenticated, AdminOnly.Error())
		}

		ctx = WithAuthUser(ctx, authUser)

		return handler(ctx, req)
	}
}
