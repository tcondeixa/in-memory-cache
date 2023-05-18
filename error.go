package cache

const (
	invalidCacheKeyError    = "key does not exist"
	outdatedCacheEntryError = "entry is expired"
	keyAlreadyExistsError   = "key already exists"
	unknownFileFormatError  = "format unknown"
)

type InvalidCacheKeyError struct{}

func (i *InvalidCacheKeyError) Error() string {
	return invalidCacheKeyError
}

type OutdatedCacheEntryError struct{}

func (i *OutdatedCacheEntryError) Error() string {
	return outdatedCacheEntryError
}

type KeyAlreadyExistsError struct{}

func (i *KeyAlreadyExistsError) Error() string {
	return keyAlreadyExistsError
}

type UnknownFileFormatError struct{}

func (i *UnknownFileFormatError) Error() string {
	return unknownFileFormatError
}
