package db

type DB interface {
	DomainDB
	ActualLRPDB
}
