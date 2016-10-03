package goca

import (
	"errors"
	"strconv"
)

type VM struct {
	XMLResource
	Id   uint
	Name string
}

type VMPool struct {
	XMLResource
}

type VM_STATE int

const (
	INIT            VM_STATE = 0
	PENDING         VM_STATE = 1
	HOLD            VM_STATE = 2
	ACTIVE          VM_STATE = 3
	STOPPED         VM_STATE = 4
	SUSPENDED       VM_STATE = 5
	DONE            VM_STATE = 6
	//FAILED        VM_STATE = 7
	POWEROFF        VM_STATE = 8
	UNDEPLOYED      VM_STATE = 9
	CLONING         VM_STATE = 10
	CLONING_FAILURE VM_STATE = 11
)

func (s VM_STATE) String() string {

	switch s {
	case INIT:            return "INIT"
	case PENDING:         return "PENDING"
	case HOLD:            return "HOLD"
	case ACTIVE:          return "ACTIVE"
	case STOPPED:         return "STOPPED"
	case SUSPENDED:       return "SUSPENDED"
	case DONE:            return "DONE"
	case POWEROFF:        return "POWEROFF"
	case UNDEPLOYED:      return "UNDEPLOYED"
	case CLONING:         return "CLONING"
	case CLONING_FAILURE: return "CLONING_FAILURE"
	default: return ""
	}
}

type LCM_STATE int

const (
	LCM_INIT                         LCM_STATE  =  0
	PROLOG                           LCM_STATE  =  1
	BOOT                             LCM_STATE  =  2
	RUNNING                          LCM_STATE  =  3
	MIGRATE                          LCM_STATE  =  4
	SAVE_STOP                        LCM_STATE  =  5
	SAVE_SUSPEND                     LCM_STATE  =  6
	SAVE_MIGRATE                     LCM_STATE  =  7
	PROLOG_MIGRATE                   LCM_STATE  =  8
	PROLOG_RESUME                    LCM_STATE  =  9
	EPILOG_STOP                      LCM_STATE  =  10
	EPILOG                           LCM_STATE  =  11
	SHUTDOWN                         LCM_STATE  =  12
	//CANCEL                         LCM_STATE  =  13
	//FAILURE                        LCM_STATE  =  14
	CLEANUP_RESUBMIT                 LCM_STATE  =  15
	UNKNOWN                          LCM_STATE  =  16
	HOTPLUG                          LCM_STATE  =  17
	SHUTDOWN_POWEROFF                LCM_STATE  =  18
	BOOT_UNKNOWN                     LCM_STATE  =  19
	BOOT_POWEROFF                    LCM_STATE  =  20
	BOOT_SUSPENDED                   LCM_STATE  =  21
	BOOT_STOPPED                     LCM_STATE  =  22
	CLEANUP_DELETE                   LCM_STATE  =  23
	HOTPLUG_SNAPSHOT                 LCM_STATE  =  24
	HOTPLUG_NIC                      LCM_STATE  =  25
	HOTPLUG_SAVEAS                   LCM_STATE  =  26
	HOTPLUG_SAVEAS_POWEROFF          LCM_STATE  =  27
	HOTPLUG_SAVEAS_SUSPENDED         LCM_STATE  =  28
	SHUTDOWN_UNDEPLOY                LCM_STATE  =  29
	EPILOG_UNDEPLOY                  LCM_STATE  =  30
	PROLOG_UNDEPLOY                  LCM_STATE  =  31
	BOOT_UNDEPLOY                    LCM_STATE  =  32
	HOTPLUG_PROLOG_POWEROFF          LCM_STATE  =  33
	HOTPLUG_EPILOG_POWEROFF          LCM_STATE  =  34
	BOOT_MIGRATE                     LCM_STATE  =  35
	BOOT_FAILURE                     LCM_STATE  =  36
	BOOT_MIGRATE_FAILURE             LCM_STATE  =  37
	PROLOG_MIGRATE_FAILURE           LCM_STATE  =  38
	PROLOG_FAILURE                   LCM_STATE  =  39
	EPILOG_FAILURE                   LCM_STATE  =  40
	EPILOG_STOP_FAILURE              LCM_STATE  =  41
	EPILOG_UNDEPLOY_FAILURE          LCM_STATE  =  42
	PROLOG_MIGRATE_POWEROFF          LCM_STATE  =  43
	PROLOG_MIGRATE_POWEROFF_FAILURE  LCM_STATE  =  44
	PROLOG_MIGRATE_SUSPEND           LCM_STATE  =  45
	PROLOG_MIGRATE_SUSPEND_FAILURE   LCM_STATE  =  46
	BOOT_UNDEPLOY_FAILURE            LCM_STATE  =  47
	BOOT_STOPPED_FAILURE             LCM_STATE  =  48
	PROLOG_RESUME_FAILURE            LCM_STATE  =  49
	PROLOG_UNDEPLOY_FAILURE          LCM_STATE  =  50
	DISK_SNAPSHOT_POWEROFF           LCM_STATE  =  51
	DISK_SNAPSHOT_REVERT_POWEROFF    LCM_STATE  =  52
	DISK_SNAPSHOT_DELETE_POWEROFF    LCM_STATE  =  53
	DISK_SNAPSHOT_SUSPENDED          LCM_STATE  =  54
	DISK_SNAPSHOT_REVERT_SUSPENDED   LCM_STATE  =  55
	DISK_SNAPSHOT_DELETE_SUSPENDED   LCM_STATE  =  56
	DISK_SNAPSHOT                    LCM_STATE  =  57
	//DISK_SNAPSHOT_REVERT           LCM_STATE  =  58
	DISK_SNAPSHOT_DELETE             LCM_STATE  =  59
	PROLOG_MIGRATE_UNKNOWN           LCM_STATE  =  60
	PROLOG_MIGRATE_UNKNOWN_FAILURE   LCM_STATE  =  61
	)

func (l LCM_STATE) String() string {
	switch l {
	case LCM_INIT:                        return "LCM_INIT"
	case PROLOG:                          return "PROLOG"
	case BOOT:                            return "BOOT"
	case RUNNING:                         return "RUNNING"
	case MIGRATE:                         return "MIGRATE"
	case SAVE_STOP:                       return "SAVE_STOP"
	case SAVE_SUSPEND:                    return "SAVE_SUSPEND"
	case SAVE_MIGRATE:                    return "SAVE_MIGRATE"
	case PROLOG_MIGRATE:                  return "PROLOG_MIGRATE"
	case PROLOG_RESUME:                   return "PROLOG_RESUME"
	case EPILOG_STOP:                     return "EPILOG_STOP"
	case EPILOG:                          return "EPILOG"
	case SHUTDOWN:                        return "SHUTDOWN"
	case CLEANUP_RESUBMIT:                return "CLEANUP_RESUBMIT"
	case UNKNOWN:                         return "UNKNOWN"
	case HOTPLUG:                         return "HOTPLUG"
	case SHUTDOWN_POWEROFF:               return "SHUTDOWN_POWEROFF"
	case BOOT_UNKNOWN:                    return "BOOT_UNKNOWN"
	case BOOT_POWEROFF:                   return "BOOT_POWEROFF"
	case BOOT_SUSPENDED:                  return "BOOT_SUSPENDED"
	case BOOT_STOPPED:                    return "BOOT_STOPPED"
	case CLEANUP_DELETE:                  return "CLEANUP_DELETE"
	case HOTPLUG_SNAPSHOT:                return "HOTPLUG_SNAPSHOT"
	case HOTPLUG_NIC:                     return "HOTPLUG_NIC"
	case HOTPLUG_SAVEAS:                  return "HOTPLUG_SAVEAS"
	case HOTPLUG_SAVEAS_POWEROFF:         return "HOTPLUG_SAVEAS_POWEROFF"
	case HOTPLUG_SAVEAS_SUSPENDED:        return "HOTPLUG_SAVEAS_SUSPENDED"
	case SHUTDOWN_UNDEPLOY:               return "SHUTDOWN_UNDEPLOY"
	case EPILOG_UNDEPLOY:                 return "EPILOG_UNDEPLOY"
	case PROLOG_UNDEPLOY:                 return "PROLOG_UNDEPLOY"
	case BOOT_UNDEPLOY:                   return "BOOT_UNDEPLOY"
	case HOTPLUG_PROLOG_POWEROFF:         return "HOTPLUG_PROLOG_POWEROFF"
	case HOTPLUG_EPILOG_POWEROFF:         return "HOTPLUG_EPILOG_POWEROFF"
	case BOOT_MIGRATE:                    return "BOOT_MIGRATE"
	case BOOT_FAILURE:                    return "BOOT_FAILURE"
	case BOOT_MIGRATE_FAILURE:            return "BOOT_MIGRATE_FAILURE"
	case PROLOG_MIGRATE_FAILURE:          return "PROLOG_MIGRATE_FAILURE"
	case PROLOG_FAILURE:                  return "PROLOG_FAILURE"
	case EPILOG_FAILURE:                  return "EPILOG_FAILURE"
	case EPILOG_STOP_FAILURE:             return "EPILOG_STOP_FAILURE"
	case EPILOG_UNDEPLOY_FAILURE:         return "EPILOG_UNDEPLOY_FAILURE"
	case PROLOG_MIGRATE_POWEROFF:         return "PROLOG_MIGRATE_POWEROFF"
	case PROLOG_MIGRATE_POWEROFF_FAILURE: return "PROLOG_MIGRATE_POWEROFF_FAILURE"
	case PROLOG_MIGRATE_SUSPEND:          return "PROLOG_MIGRATE_SUSPEND"
	case PROLOG_MIGRATE_SUSPEND_FAILURE:  return "PROLOG_MIGRATE_SUSPEND_FAILURE"
	case BOOT_UNDEPLOY_FAILURE:           return "BOOT_UNDEPLOY_FAILURE"
	case BOOT_STOPPED_FAILURE:            return "BOOT_STOPPED_FAILURE"
	case PROLOG_RESUME_FAILURE:           return "PROLOG_RESUME_FAILURE"
	case PROLOG_UNDEPLOY_FAILURE:         return "PROLOG_UNDEPLOY_FAILURE"
	case DISK_SNAPSHOT_POWEROFF:          return "DISK_SNAPSHOT_POWEROFF"
	case DISK_SNAPSHOT_REVERT_POWEROFF:   return "DISK_SNAPSHOT_REVERT_POWEROFF"
	case DISK_SNAPSHOT_DELETE_POWEROFF:   return "DISK_SNAPSHOT_DELETE_POWEROFF"
	case DISK_SNAPSHOT_SUSPENDED:         return "DISK_SNAPSHOT_SUSPENDED"
	case DISK_SNAPSHOT_REVERT_SUSPENDED:  return "DISK_SNAPSHOT_REVERT_SUSPENDED"
	case DISK_SNAPSHOT_DELETE_SUSPENDED:  return "DISK_SNAPSHOT_DELETE_SUSPENDED"
	case DISK_SNAPSHOT:                   return "DISK_SNAPSHOT"
	case DISK_SNAPSHOT_DELETE:            return "DISK_SNAPSHOT_DELETE"
	case PROLOG_MIGRATE_UNKNOWN:          return "PROLOG_MIGRATE_UNKNOWN"
	case PROLOG_MIGRATE_UNKNOWN_FAILURE:  return "PROLOG_MIGRATE_UNKNOWN_FAILURE"
	default: return ""
	}
}

func NewVMPool(args ...int) (*VMPool, error) {
	var who, start_id, end_id, state int

	switch len(args) {
	case 0:
		who = PoolWhoMine
		start_id = -1
		end_id = -1
		state = -1
	case 1:
		who = args[0]
		start_id = -1
		end_id = -1
		state = -1
	case 3:
		who = args[0]
		start_id = args[1]
		end_id = args[2]
		state = -1
	case 4:
		who = args[0]
		start_id = args[1]
		end_id = args[2]
		state = args[3]
	default:
		return nil, errors.New("Wrong number of arguments")
	}

	response, err := client.Call("one.vmpool.info", who, start_id, end_id, state)
	if err != nil {
		return nil, err
	}

	vmpool := &VMPool{XMLResource{body: response.Body()}}

	return vmpool, err

}

func CreateVM(template string, pending bool) (uint, error) {
	response, err := client.Call("one.vm.allocate", template, pending)
	if err != nil {
		return 0, err
	}

	return uint(response.BodyInt()), nil
}

func NewVM(id uint) *VM {
	return &VM{Id: id}
}

func NewVMFromName(name string) (*VM, error) {
	vmpool, err := NewVMPool()
	if err != nil {
		return nil, err
	}

	id, err := vmpool.GetIdFromName(name, "/VM_POOL/VM")
	if err != nil {
		return nil, err
	}

	return NewVM(id), nil
}

func (vm *VM) Info() error {
	response, err := client.Call("one.vm.info", vm.Id)
	vm.body = response.Body()
	return err
}

func (vm *VM) State() (int, int, error) {
	vm_stateString, ok := vm.XPath("/VM/STATE")
	if ok != true {
		return -1, -1, errors.New("Unable to parse VM State")
	}

	lcm_stateString, ok := vm.XPath("/VM/LCM_STATE")
	if ok != true {
		return -1, -1, errors.New("Unable to parse LCM State")
	}

	vm_state, _ := strconv.Atoi(vm_stateString)
	lcm_state, _ := strconv.Atoi(lcm_stateString)

	return vm_state, lcm_state, nil
}

func (vm *VM) StateString() (string, string, error) {
	vm_state, lcm_state, err := vm.State()
	if err != nil {
		return "", "", err
	}
	return VM_STATE(vm_state).String(), LCM_STATE(lcm_state).String(), nil
}

func (vm *VM) Action(action string) error {
	_, err := client.Call("one.vm.action", action, vm.Id)
	return err
}

// VM Actions

func (vm *VM) TerminateHard() error {
	return vm.Action("terminate-hard")
}

func (vm *VM) Terminate() error {
	return vm.Action("terminate")
}

func (vm *VM) UndeployHard() error {
	return vm.Action("undeploy-hard")
}

func (vm *VM) Undeploy() error {
	return vm.Action("undeploy")
}

func (vm *VM) PoweroffHard() error {
	return vm.Action("poweroff-hard")
}

func (vm *VM) Poweroff() error {
	return vm.Action("poweroff")
}

func (vm *VM) RebootHard() error {
	return vm.Action("reboot-hard")
}

func (vm *VM) Reboot() error {
	return vm.Action("reboot")
}

func (vm *VM) Hold() error {
	return vm.Action("hold")
}

func (vm *VM) Release() error {
	return vm.Action("release")
}

func (vm *VM) Stop() error {
	return vm.Action("stop")
}

func (vm *VM) Suspend() error {
	return vm.Action("suspend")
}

func (vm *VM) Resume() error {
	return vm.Action("resume")
}

func (vm *VM) Resched() error {
	return vm.Action("resched")
}

func (vm *VM) Unresched() error {
	return vm.Action("unresched")
}
