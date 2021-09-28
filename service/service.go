package service

type Service interface {
	// Start a service
	Start() error

	// Shutdown a service
	Shutdown()
}
