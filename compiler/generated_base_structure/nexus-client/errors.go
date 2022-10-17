package nexus_client

import (
	"errors"
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IsNotFound(err error) bool {
	return k8serrors.IsNotFound(err)
}

func IsUnauthorized(err error) bool {
	return k8serrors.IsUnauthorized(err)
}

func IsAlreadyExists(err error) bool {
	return k8serrors.IsAlreadyExists(err)
}

func IsConflict(err error) bool {
	return k8serrors.IsConflict(err)
}

func IsInvalid(err error) bool {
	return k8serrors.IsInvalid(err)
}

func IsGone(err error) bool {
	return k8serrors.IsGone(err)
}

func IsResourceExpired(err error) bool {
	return k8serrors.IsResourceExpired(err)
}

func IsNotAcceptable(err error) bool {
	return k8serrors.IsNotAcceptable(err)
}

func IsUnsupportedMediaType(err error) bool {
	return k8serrors.IsUnsupportedMediaType(err)
}

func IsMethodNotSupported(err error) bool {
	return k8serrors.IsMethodNotSupported(err)
}

func IsServiceUnavailable(err error) bool {
	return k8serrors.IsServiceUnavailable(err)
}

func IsBadRequest(err error) bool {
	return k8serrors.IsBadRequest(err)
}

func IsForbidden(err error) bool {
	return k8serrors.IsForbidden(err)
}

func IsTimeout(err error) bool {
	return k8serrors.IsTimeout(err)
}

func IsServerTimeout(err error) bool {
	return k8serrors.IsServerTimeout(err)
}

func IsInternalError(err error) bool {
	return k8serrors.IsInternalError(err)
}

func IsTooManyRequests(err error) bool {
	return k8serrors.IsTooManyRequests(err)
}

func IsRequestEntityTooLargeError(err error) bool {
	return k8serrors.IsRequestEntityTooLargeError(err)
}

func IsUnexpectedServerError(err error) bool {
	return k8serrors.IsUnexpectedServerError(err)
}

func IsUnexpectedObjectError(err error) bool {
	return k8serrors.IsUnexpectedObjectError(err)
}

type ParentNotFound struct {
	errMessage string
}

func NewParentNotFound(displayName string, objectType metav1.Type) ParentNotFound {
	return ParentNotFound{
		errMessage: fmt.Sprintf("parent not found for %s: %s", objectType, displayName),
	}
}

func (p ParentNotFound) Error() string {
	return p.errMessage
}

func IsParentNotFound(err error) bool {
	return errors.As(err, &ParentNotFound{})
}

type ChildNotFound struct {
	errMessage string
}

func NewChildNotFound(parentDisplayName string, parentType string,
	childVarName string, childDisplayName ...string) ChildNotFound {
	if len(childDisplayName) == 0 {
		return ChildNotFound{
			errMessage: fmt.Sprintf("child %s not found for %s: %s", childVarName, parentType, parentDisplayName),
		}
	}
	return ChildNotFound{
		errMessage: fmt.Sprintf("child %s: %s not found for %s: %s",
			childVarName, childDisplayName[0], parentType, parentDisplayName),
	}
}

func (p ChildNotFound) Error() string {
	return p.errMessage
}

func IsChildNotFound(err error) bool {
	return errors.As(err, &ChildNotFound{})
}

type LinkNotFound struct {
	errMessage string
}

func NewLinkNotFound(parentDisplayName string, parentType string,
	linkVarName string, linkDisplayName ...string) LinkNotFound {
	if len(linkDisplayName) == 0 {
		return LinkNotFound{
			errMessage: fmt.Sprintf("link %s not found for %s: %s", linkVarName, parentType, parentDisplayName),
		}
	}
	return LinkNotFound{
		errMessage: fmt.Sprintf("link %s: %s not found for %s: %s",
			linkVarName, linkDisplayName[0], parentType, parentDisplayName),
	}
}

func (p LinkNotFound) Error() string {
	return p.errMessage
}

func IsLinkNotFound(err error) bool {
	return errors.As(err, &LinkNotFound{})
}

type SingletonNameError struct {
	errMessage string
}

func NewSingletonNameError(displayName string) SingletonNameError {
	return SingletonNameError{
		errMessage: fmt.Sprintf("wrong name of singleton object: %s, singleton can have only"+
			"'default' as a display name", displayName),
	}
}

func (p SingletonNameError) Error() string {
	return p.errMessage
}

func IsSingletonNameError(err error) bool {
	return errors.As(err, &SingletonNameError{})
}
