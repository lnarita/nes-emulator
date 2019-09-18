class OpCode(object):
    def __init__(self, code, addressing_mode, cycles):
        self.id = code
        self.addressing_mode = addressing_mode
        self.cycles = cycles

    @classmethod
    def exec(cls, cpu, memory):
        raise NotImplementedError("Class {} not implemented yet! CPU: {}, MEM: {}".format(cls, cpu, memory))

    @classmethod
    def create_variations(cls):
        pass

    @classmethod
    def create_dict_entry(cls, x):
        return tuple((x[0], cls(*x)))
