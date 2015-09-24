package goca

type VM struct {
	Id   uint
	body XML
}

func (vm *VM) Info() error {
	response, err := client.Call("one.vm.info", vm.Id)

	vm.body = XML(response.Body)

	return err
}

func (vm *VM) Body() XML {
	return vm.body
}
