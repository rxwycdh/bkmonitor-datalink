package window

// Window expiration is implemented using the observer mechanism,
// and the logic of different Window needs to comply with observer / trigger interface
// when the processing time is triggered
type (
	Event struct {
		*CollectTrace
	}

	observer interface {
		detectNotify()
		handleNotify()
		assembleRuntimeConfig(runtimeOpt ...RuntimeConfigOption)
	}

	trigger interface {
		startWatch()

		register(observer)
		deRegister(observer)
	}
)
