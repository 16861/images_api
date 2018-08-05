from PIL import Image
from os import listdir

PATH = "imgs/"

for pth in listdir(PATH):
    img = Image.open(PATH + pth)
    img.save("imgs/"+"COMPRESSED_"+pth.split("//")[-1], quality=30)

