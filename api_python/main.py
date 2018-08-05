from PIL import Image
from io import BytesIO
from flask import Flask, request, jsonify
import base64 as b64

app = Flask(__name__)

def compressImg(im):
    with BytesIO() as f:
        im.save(f, format="jpeg", quality=65)
        f.seek(0)
        return f.read()
        """ with open("tmp.jpg", 'bw') as fd:
            fd.write(f.read()) """

@app.route("/compressImage", methods=["POST"])
def compressImage():
    js = request.get_json()
    im = Image.open(BytesIO(b64.b64decode(js["image"])))
    ret = compressImg(im)
    js["image"] = b64.b64encode(ret)
    js["filename"] = js["filename"].split(".")[0] + ".jpg"
    return jsonify(filename = js["filename"], image = js["image"].decode("utf-8"))


@app.route("/compressImages", methods=["POST"])
def compressImages():
    js = request.get_json()
    for j in js:
        im = Image.open(BytesIO(b64.b64decode(j["image"])))
        ret = compressImg(im)
        j["image"] = b64.b64encode(ret).decode("utf-8")
        j["filename"] = j["filename"].split(".")[0] + ".jpg"
    return jsonify(js)



app.run(debug=True)