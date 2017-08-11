// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

type hashMismatchError struct {
	emsg string
}

type alreadyRunningError struct {
	emsg string
}

type authFailedError struct {
	emsg string
}

func (e *hashMismatchError) Error() string {
	return e.emsg
}

func (e *alreadyRunningError) Error() string {
	return e.emsg
}

func (e *authFailedError) Error() string {
	return e.emsg
}

func IsErrAlreadyRunning(err error) bool {
	if _, ok := err.(*alreadyRunningError); ok {
		return true
	}
	return false
}

func IsErrHashMismatch(err error) bool {
	if _, ok := err.(*hashMismatchError); ok {
		return true
	}
	return false
}

func IsErrAuthFailed(err error) bool {
	if _, ok := err.(*authFailedError); ok {
		return true
	}
	return false
}
