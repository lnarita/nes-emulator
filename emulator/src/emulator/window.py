# pygame silent import
import os, sys
with open(os.devnull, 'w') as f:
    oldstdout = sys.stdout
    sys.stdout = f

    import pygame

    sys.stdout = oldstdout

from pygame import Color
import time
from emulator.config import *

W_WIDTH = 256*DISPLAY_SCALE
W_HEIGHT = 240*DISPLAY_SCALE


class Window():
	def __init__(self):
		check_errors = pygame.init()
		if check_errors[1] > 0:
		    print("pygame error")
		    sys.exit(-1)

		size = (W_WIDTH, W_HEIGHT)

		self.surface = pygame.display.set_mode(size)
		pygame.display.set_caption("NES emulator")

	def setPixel(self,i,j,color):
		pygame.draw.rect(self.surface, color, pygame.Rect(i*DISPLAY_SCALE,j*DISPLAY_SCALE,DISPLAY_SCALE,DISPLAY_SCALE))


	def flip(self):
		pygame.display.flip()


def main():
	screen = Window()
	for i in range(262):
		screen.scanLine(i)
	time.sleep(10)

if __name__ == "__main__":
	main()
