from PIL import Image
from os import listdir

PATH = "/home/igor/go/src/rest-and-go/imgs/"

for pth in listdir(PATH):
    img = Image.open(PATH + pth)
    img.save("/home/igor/go/src/rest-and-go/imgs/"+"COMPRESSED_"+pth.split("//")[-1], quality=30)

