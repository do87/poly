package usecases

// Usecase service
type Usecase struct {
	Agents *agentsUsecase
	Keys   *keysUsecase
}
