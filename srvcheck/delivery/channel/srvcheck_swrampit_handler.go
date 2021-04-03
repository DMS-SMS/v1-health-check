// in srvcheck_swarmpit_handler.go file, define delivery from channel msg to swarmpit check usecase handler
// publishing msg to golang channel which is received from outside is not occurred in this package

package channel

// swarmpitCheckHandler is delivered data handler about swarmpit check using usecase layer
type swarmpitCheckHandler struct {
	// SUsecase is usecase layer interface which is injected from package outside (maybe, in main)
	SUsecase domain.SwarmpitCheckUseCase
}
