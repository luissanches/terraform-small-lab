from array import array
from flask import Flask, jsonify, request
import uuid, randomname


app = Flask(__name__)

@app.get('/uuid')
def generateUuid():
	id = request.args.get('id')
	response = {
		'status': 'ok',
		'uuid': uuid.uuid4().hex,
		'instanceId': id,
	}

	return jsonify(response)

@app.get('/instances')
def instances():
	userName = request.args.get('username')
	response = [{
		'status': 'ok',
		'uuid': uuid.uuid4().hex,
		'username': userName,
		'owner': randomname.get_name()
	},
	{
		'status': 'ok',
		'uuid': uuid.uuid4().hex,
		'username': userName,
		'owner': randomname.get_name()
	}]

	return jsonify(response)

@app.get('/owner')
def getOwner():
	response = {
		'status': 'ok',
		'name': randomname.get_name()
	}

	return jsonify(response)

@app.post('/instance')
def instance():
	json = request.json
	print(json)

	response = {
		'status': 'ok',
		'type': json['type'],
		'id': uuid.uuid4().hex,
		'resource_name': randomname.get_name()
	}

	return jsonify(response)