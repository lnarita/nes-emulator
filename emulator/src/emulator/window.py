# pygame silent import
import os, sys
with open(os.devnull, 'w') as f:
    oldstdout = sys.stdout
    sys.stdout = f

    import pygame

    sys.stdout = oldstdout

from pygame import Color
import time

W_WIDTH = 256*2
W_HEIGHT = 240*2




class Window():
    def __init__(self, read_count=0, strobe=False):
        check_errors = pygame.init()
        if check_errors[1] > 0:
            print("pygame error")
            sys.exit(-1)

        size = (W_WIDTH, W_HEIGHT)

        self.surface = pygame.display.set_mode(size)
        pygame.display.set_caption("NES emulator")

        self.keys = [False] * 8
        self.read_count = read_count
        self.strobe = strobe

    def setPixel(self,i,j,color):
        pygame.draw.rect(self.surface, color, pygame.Rect(i*2,j*2,2,2))


    def flip(self):
        pygame.display.flip()


    def read_keys(self):
        # print(self.read_count)
        if self.read_count >= 8:
            return 1

        if not self.strobe:
            key = self.keys[self.read_count]
            if sum(self.keys) != 0:
                print("keys", self.keys)
            return int(key)
        else:
            pygame.event.pump()
            pressed = pygame.key.get_pressed()
            # print("pressed", pressed)

            self.keys[0] = pressed[pygame.K_3]      # A
            self.keys[1] = pressed[pygame.K_2]      # B
            self.keys[2] = pressed[pygame.K_TAB]    # Select
            self.keys[3] = pressed[pygame.K_RETURN] # Start
            self.keys[4] = pressed[pygame.K_w]      # Up
            self.keys[5] = pressed[pygame.K_s]      # Down
            self.keys[6] = pressed[pygame.K_a]      # Left
            self.keys[7] = pressed[pygame.K_d]      # Right

            key = self.keys[self.read_count]
            if sum(self.keys) != 0:
                print("keys", self.keys)
            return int(key)


    def latch_keys(self, value):
        print("latch", value)
        self.strobe = value
        self.read_count = 0
        pygame.event.pump()
        pressed = pygame.key.get_pressed()
        # print("pressed", pressed)
        self.keys[0] = pressed[pygame.K_3]      # A
        self.keys[1] = pressed[pygame.K_2]      # B
        self.keys[2] = pressed[pygame.K_TAB]    # Select
        self.keys[3] = pressed[pygame.K_RETURN] # Start
        self.keys[4] = pressed[pygame.K_w]      # Up
        self.keys[5] = pressed[pygame.K_s]      # Down
        self.keys[6] = pressed[pygame.K_a]      # Left
        self.keys[7] = pressed[pygame.K_d]      # Right

def main():
    screen = Window()
    for i in range(262):
        screen.scanLine(i)
    time.sleep(10)

if __name__ == "__main__":
    main()
