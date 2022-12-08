import os
from flask import Flask, request, jsonify
from flask_cors import CORS, cross_origin
app = Flask(__name__)
cors = CORS(app)
app.config['CORS_HEADERS'] = 'Content-Type'


@app.route('/runPanda', methods=["POST"])
@cross_origin()
def runPanda():
    params = request.get_json(force=True)
    print(params)
    command = "docker run -it -v " + params['image'] + " pandare/panda " + params['commands']
    os.system(command)

    return jsonify(message=command)