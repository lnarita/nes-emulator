from emulator.cpu import CPU
from emulator.memory import Memory
from emulator.opcodes.opcodes import OpCodes


def test_ROL_accumulator():
    instruction = OpCodes.all[0x2A]
    cpu = CPU()
    memory = Memory()
    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b00000000
    cpu.carry = False
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00000000
    assert cpu.zero == True
    assert cpu.carry == False
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b00000001
    cpu.carry = False
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00000010
    assert cpu.zero == False
    assert cpu.carry == False
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b01000000
    cpu.carry = False
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b10000000
    assert cpu.zero == False
    assert cpu.carry == False
    assert cpu.negative == True

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b10000000
    cpu.carry = False
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00000000
    assert cpu.zero == True
    assert cpu.carry == True
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b00000000
    cpu.carry = True
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00000001
    assert cpu.zero == False
    assert cpu.carry == False
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b00000010
    cpu.carry = True
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00000101
    assert cpu.zero == False
    assert cpu.carry == False
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b01000000
    cpu.carry = True
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b10000001
    assert cpu.zero == False
    assert cpu.carry == False
    assert cpu.negative == True

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b10000000
    cpu.carry = True
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00000001
    assert cpu.zero == False
    assert cpu.carry == True
    assert cpu.negative == False


def test_ROR_accumulator():
    instruction = OpCodes.all[0x6A]
    cpu = CPU()
    memory = Memory()
    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b00000000
    cpu.carry = False
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00000000
    assert cpu.zero == True
    assert cpu.carry == False
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b00000001
    cpu.carry = False
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00000000
    assert cpu.zero == True
    assert cpu.carry == True
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b01000000
    cpu.carry = False
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00100000
    assert cpu.zero == False
    assert cpu.carry == False
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b10000000
    cpu.carry = False
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b01000000
    assert cpu.zero == False
    assert cpu.carry == False
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b00000000
    cpu.carry = True
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b10000000
    assert cpu.zero == False
    assert cpu.carry == False
    assert cpu.negative == True

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b00000010
    cpu.carry = True
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b10000001
    assert cpu.zero == False
    assert cpu.carry == False
    assert cpu.negative == True

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b01000000
    cpu.carry = True
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b10100000
    assert cpu.zero == False
    assert cpu.carry == False
    assert cpu.negative == True

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b00000001
    cpu.carry = True
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b10000000
    assert cpu.zero == False
    assert cpu.carry == True
    assert cpu.negative == True


def test_LSR_accumulator():
    instruction = OpCodes.all[0x4A]
    cpu = CPU()
    memory = Memory()
    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b00000000
    cpu.carry = False
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00000000
    assert cpu.zero == True
    assert cpu.carry == False
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b00000001
    cpu.carry = False
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00000000
    assert cpu.zero == True
    assert cpu.carry == True
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b01000000
    cpu.carry = False
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00100000
    assert cpu.zero == False
    assert cpu.carry == False
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b10000000
    cpu.carry = False
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b01000000
    assert cpu.zero == False
    assert cpu.carry == False
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b00000000
    cpu.carry = True
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00000000
    assert cpu.zero == True
    assert cpu.carry == False
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b00000010
    cpu.carry = True
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00000001
    assert cpu.zero == False
    assert cpu.carry == False
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b01000000
    cpu.carry = True
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00100000
    assert cpu.zero == False
    assert cpu.carry == False
    assert cpu.negative == False

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.a = 0b00000001
    cpu.carry = True
    instruction.exec(cpu, memory)

    assert cpu.cycle == (instruction.cycles - 1)
    assert cpu.a == 0b00000000
    assert cpu.zero == True
    assert cpu.carry == True
    assert cpu.negative == False