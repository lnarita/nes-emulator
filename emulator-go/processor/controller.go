package processor

type Controller struct {
	buttons [8]bool
	current byte // 0 - A; 1 - B; 2 - Select; 3 - Start; 4 - Up; 5 - Down; 6 - Left; 7 - Right
	strobe  byte
}

func (controller *Controller) SetButtons(buttons [8]bool) {
	controller.buttons = buttons
}

func (controller *Controller) Read() byte {
	var value byte = 0
	if controller.current < 8 {
		if controller.buttons[controller.current] {
			value = 1
		}
	} else {
		value = 1
	}
	controller.current++
	controller.checkStrobe()
	return value
}

func (controller *Controller) Write(data byte) {
	controller.strobe = data
	controller.checkStrobe()
}

func (controller *Controller) checkStrobe() {
	if controller.strobe&0x0001 == 0x0001 {
		controller.current = 0
	}
}
