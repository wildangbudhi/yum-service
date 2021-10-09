package domain

type RepositoryErrorType int

const (
	RepositoryError            RepositoryErrorType = 1 << iota
	RepositoryDataNotFound     RepositoryErrorType = 1 << iota
	RepositoryCreateDataFailed RepositoryErrorType = 1 << iota
	RepositoryUpdateDataFailed RepositoryErrorType = 1 << iota
)

type HTTPStatusCode int

const (
	OTPCustomertype int = 1
	OTPRestotype    int = 2
)
