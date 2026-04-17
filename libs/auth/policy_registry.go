package grpcAuth

import (
	"fmt"
	"log"

	commonv1 "github.com/deniSSTK/task-engine/gen/proto/common/v1"
	defErrors "github.com/deniSSTK/task-engine/libs/errors"
	grpcErr "github.com/deniSSTK/task-engine/libs/grpc/error"
	"github.com/deniSSTK/task-engine/libs/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

type PolicyRegistry struct {
	methods map[string]commonv1.AuthPolicy

	log          *logger.Logger
	errorWrapper *grpcErr.AppErrorWrapper
}

func NewPolicyRegistry(
	log *logger.Logger,
	errorWrapper *grpcErr.AppErrorWrapper,
	files ...protoreflect.FileDescriptor) *PolicyRegistry {
	log = log.Named("PolicyRegistry")

	r := &PolicyRegistry{
		methods:      make(map[string]commonv1.AuthPolicy),
		log:          log,
		errorWrapper: errorWrapper,
	}

	for _, file := range files {
		services := file.Services()

		for i := 0; i < services.Len(); i++ {
			service := services.Get(i)
			methods := service.Methods()

			for j := 0; j < methods.Len(); j++ {
				method := methods.Get(j)

				policy, err := extractPolicy(method)
				if err != nil {
					log.Fatal(
						err.Error(),
						zap.String("service", string(service.FullName())),
						zap.String("method", string(method.Name())),
					)
				}

				fullMethod := fmt.Sprintf("/%s/%s", service.FullName(), method.Name())
				r.methods[fullMethod] = policy
			}
		}
	}

	return r
}

func (r *PolicyRegistry) MustGet(fullName string) commonv1.AuthPolicy {
	policy, ok := r.methods[fullName]
	if !ok {
		log.Fatal(
			defErrors.MissingRequiredOption.Error(),
			zap.String("fullName", fullName),
		)
	}

	return policy
}

func extractPolicy(method protoreflect.MethodDescriptor) (commonv1.AuthPolicy, error) {
	opts, ok := method.Options().(*descriptorpb.MethodOptions)
	if !ok || opts == nil {
		return commonv1.AuthPolicy_AUTH_POLICY_UNSPECIFIED, defErrors.MethodOptionsNotFound
	}

	if !proto.HasExtension(opts, commonv1.E_AuthPolicy) {
		return commonv1.AuthPolicy_AUTH_POLICY_UNSPECIFIED, defErrors.MissingRequiredOption
	}

	ext := proto.GetExtension(opts, commonv1.E_AuthPolicy)

	policy, ok := ext.(commonv1.AuthPolicy)
	if !ok {
		return commonv1.AuthPolicy_AUTH_POLICY_UNSPECIFIED, defErrors.InvalidExtensionType
	}

	if policy == commonv1.AuthPolicy_AUTH_POLICY_UNSPECIFIED {
		return commonv1.AuthPolicy_AUTH_POLICY_UNSPECIFIED, defErrors.MetadataOptionUnspecified
	}

	return policy, nil
}
