from emulator.cpu import CPU
from emulator.memory import Memory
from emulator.opcodes.opcodes import OpCodes


def test_DEX_within_bounds():
    instruction = OpCodes.all[0xCA]
    cpu = CPU()
    memory = Memory()
    cpu.inc_cycle_by(-cpu.cycle)
    cpu.x = 7
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.x == 6
    assert cpu.zero == False
    assert cpu.negative == False


def test_DEX_negative():
    instruction = OpCodes.all[0xCA]
    cpu = CPU()
    memory = Memory()
    cpu.inc_cycle_by(-cpu.cycle)
    cpu.x = 0
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.x == 255
    assert cpu.zero == False
    assert cpu.negative == True

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.x = -1
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.x == 254
    assert cpu.zero == False
    assert cpu.negative == True


def test_DEX_zero():
    instruction = OpCodes.all[0xCA]
    cpu = CPU()
    memory = Memory()
    cpu.inc_cycle_by(-cpu.cycle)
    cpu.x = 1
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.x == 0
    assert cpu.zero == True
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.x = 257
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.x == 0
    assert cpu.zero == True
    assert cpu.negative == False


def test_DEY_within_bounds():
    instruction = OpCodes.all[0x88]
    cpu = CPU()
    memory = Memory()
    cpu.inc_cycle_by(-cpu.cycle)
    cpu.y = 7
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.y == 6
    assert cpu.zero == False
    assert cpu.negative == False


def test_DEY_negative():
    instruction = OpCodes.all[0x88]
    cpu = CPU()
    memory = Memory()
    cpu.inc_cycle_by(-cpu.cycle)
    cpu.y = 0
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.y == 255
    assert cpu.zero == False
    assert cpu.negative == True

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.y = -1
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.y == 254
    assert cpu.zero == False
    assert cpu.negative == True

def test_DEY_zero():
    instruction = OpCodes.all[0x88]
    cpu = CPU()
    memory = Memory()
    cpu.inc_cycle_by(-cpu.cycle)
    cpu.y = 1
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.y == 0
    assert cpu.zero == True
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.y = 257
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.y == 0
    assert cpu.zero == True
    assert cpu.negative == False

def test_INX_within_bounds():
    instruction = OpCodes.all[0xE8]
    cpu = CPU()
    memory = Memory()
    cpu.inc_cycle_by(-cpu.cycle)
    cpu.x = 7
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.x == 8
    assert cpu.zero == False
    assert cpu.negative == False


def test_INX_zero():
    instruction = OpCodes.all[0xE8]
    cpu = CPU()
    memory = Memory()
    cpu.inc_cycle_by(-cpu.cycle)
    cpu.x = -1
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.x == 0
    assert cpu.zero == True
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.x = 255
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.x == 0
    assert cpu.zero == True
    assert cpu.negative == False


def test_INY_within_bounds():
    instruction = OpCodes.all[0xC8]
    cpu = CPU()
    memory = Memory()
    cpu.inc_cycle_by(-cpu.cycle)
    cpu.y = 7
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.y == 8
    assert cpu.zero == False
    assert cpu.negative == False


def test_INY_zero():
    instruction = OpCodes.all[0xC8]
    cpu = CPU()
    memory = Memory()
    cpu.inc_cycle_by(-cpu.cycle)
    cpu.y = -1
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.y == 0
    assert cpu.zero == True
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.y = 255
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.y == 0
    assert cpu.zero == True
    assert cpu.negative == False
