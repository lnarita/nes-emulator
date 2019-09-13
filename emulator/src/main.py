import argparse

from cartridge import Cartridge
from memory import Memory
from opcodes.opcodes import OpCodes
from states import BaseCPUState


def main(args):
    emulate(args.file)


def emulate(file_path):
    running = True
    file_contents = read_file(file_path)
    cartridge = Cartridge.from_bytes(file_contents)
    memory = Memory(cartridge.prg_rom)
    cpu_state = BaseCPUState(pc=0xC000)

    while (running):
        try:
            instruction = memory.fetch(cpu_state.pc)
            decoded = decode_instruction(instruction)
            if decoded:
                print(decoded)
                # decoded.exec(cpu_state, memory)
            cpu_state.pc += 1
        except IndexError:
            running = False


def read_file(file_path):
    with open(file_path, mode='rb') as file:
        return file.read()


def decode_instruction(instruction):
    decoded = OpCodes.all[instruction]
    if not decoded:
        print("Can't find instruction 0x{:02X}".format(instruction))
    return decoded


def print_debug_line(cpu_state):
    print(cpu_state)


if __name__ == "__main__":
    """ This is executed when run from the command line """
    parser = argparse.ArgumentParser()

    # Required positional argument
    parser.add_argument("file", help="NES Cartridge file path")

    args = parser.parse_args()
    main(args)
