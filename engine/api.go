package engine


func CreateBotInstance() *BotEngine {
	return &BotEngine{
		context: context{
			state: IDLE_STATE,
		},
	}
}

func (e *BotEngine) Init() []byte {

	//initialize the engine and return the welcome message.
	e.context.api.Init("891b596c8335d66d4c3c11692eb88be1")
	return generateStartupMessage()
}

func (e *BotEngine) ProcessInput(s string) []byte {

	//send the input to engine so it decides which intent to generate
	return e.compute(s)
}
