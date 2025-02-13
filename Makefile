.PHONY: launch_emulator

launch_emulator:
ifndef HOST
	$(error "HOST is not set. Usage: ex. make HOST=localhost PORT=8030 launch_emulator")
endif
ifndef PORT
	$(error "PORT is not set. Usage: ex. make HOST=localhost PORT=8030 launch_emulator")
endif
	gcloud emulators firestore start --quiet --host-port=$(HOST):$(PORT)
