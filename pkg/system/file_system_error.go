package system

type FileSystemError struct {
	message      string
	causingError error
}

func (fse *FileSystemError) Error() string {
	var causingErrorMessage string
	if fse.causingError == nil {
		causingErrorMessage = "<nil>"
	} else {
		causingErrorMessage = fse.causingError.Error()
	}
	if fse == nil {
		return "<nil>: " + causingErrorMessage
	}
	return fse.message + ": " + causingErrorMessage
}

func (fse *FileSystemError) Unwrap() error {
	return fse.causingError
}

func NewFileSystemError(message string, cause error) *FileSystemError {
	return &FileSystemError{message: message, causingError: cause}
}
