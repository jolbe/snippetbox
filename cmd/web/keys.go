package main

type key string

// context keys
const isAuthenticatedContextKey key = "isAuthenticatedUser"

// session keys
const (
	authenticatedUserIDSessionKey key = "authenticatedUserID"
	flashSessionKey               key = "flash"
)
