import argparse

from emulator.adressing import AddressMode
from emulator.cartridge import Cartridge
from emulator.constants import HIGH_BITS_MASK, LOW_BITS_MASK
from emulator.cpu import CPU, CPUState, StatusRegisterFlags
from emulator.memory import Memory, MemoryPositions
from emulator.opcodes.jump import BRK
from emulator.opcodes.opcodes import OpCodes


def main(args):
    emulate(args.file)


# https://stackoverflow.com/questions/45305891/6502-cycle-timing-per-instruction
# http://nesdev.com/6502_cpu.txt (6510 Instruction Timing)
def emulate(file_path):
    nestest_log_format = args.nestest
    automation_mode = args.automation
    running = True
    file_contents = read_file(file_path)
    cartridge = Cartridge.from_bytes(file_contents)
    memory = Memory(cartridge.prg_rom)
    cpu = CPU(log_compatible_mode=nestest_log_format)

    reset_pos_low = memory.fetch(MemoryPositions.RESET.start)
    reset_pos_high = memory.fetch(MemoryPositions.RESET.end)
    reset_pos = MemoryPositions.PRG_ROM_START.wrap(
        AddressMode.get_16_bits_addr_from_high_low((reset_pos_high << 8) & HIGH_BITS_MASK, reset_pos_low & LOW_BITS_MASK))
    cpu.pc = reset_pos

    if nestest_log_format:
        # hack for Nintendulator nestest log comparison
        cpu.flags = 0x24
        cpu.inc_cycle_by(7)
    if automation_mode:
        cpu.pc = MemoryPositions.PRG_ROM_START.start

    while running:
        try:
            previous_state = CPUState(pc=cpu.pc, sp=cpu.sp, a=cpu.a, x=cpu.x, y=cpu.y, p=StatusRegisterFlags(int_value=cpu.flags), addr=cpu.addr, data=cpu.data,
                                      cycle=cpu.cycle, log_compatible_mode=nestest_log_format)
            decoded = cpu.exec_in_cycle(fetch_and_decode_instruction, cpu, memory)  # fetching and decoding a instruction always take 1 cycle
            if decoded:
                decoded.exec(cpu, memory)
                if isinstance(decoded, BRK):
                    # abort program on BRK
                    running = False
                    break
                print_debug_line(cpu, previous_state, decoded, nestest_log_format)
                cpu.clear_state_mem()
        except IndexError as e:
            # we've reached a program counter that is not within memory bounds
            print(e)
            running = False
        except KeyError as e:
            # print("Can't find instruction by code {}".format(e))
            cpu.inc_pc_by(1)
        except Exception as e:
            print(e)
            cpu.inc_pc_by(1)


def fetch_and_decode_instruction(cpu, memory):
    instruction = memory.fetch(cpu.pc)
    decoded = decode_instruction(instruction)
    cpu.inc_pc_by(1)
    return decoded


def read_file(file_path):
    with open(file_path, mode='rb') as file:
        return file.read()


def decode_instruction(instruction):
    decoded = OpCodes.all[instruction]
    if not decoded:
        print("Can't find instruction 0x{:02X}".format(instruction))
    return decoded


def print_debug_line(cpu, previous, instruction, nestest):
    if nestest:
        print("%04X  %s  %s  CYC:%d" % (previous.pc, instruction, previous, previous.cycle))
        # print("%04X  %s  %s  CYC:%d (Δ = %d; expected = %d)" % (previous.pc, instruction, previous, previous.cycle, (cpu.cycle - previous.cycle), instruction.cycles))
    else:
        print(cpu)


if __name__ == "__main__":
    """ This is executed when run from the command line """
    parser = argparse.ArgumentParser()

    # Required positional argument
    parser.add_argument("file", help="NES Cartridge file path")
    parser.add_argument("--nestest", action="store_true")
    parser.add_argument("--automation", action="store_true")

    args = parser.parse_args()
    main(args)
