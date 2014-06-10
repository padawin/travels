#!/usr/bin/env python

import sys
import os
import json
from subprocess import call

def getFiles(path):
	fileList = []
	rootdir = sys.argv[1]
	for root, subFolders, files in os.walk(path):
		for file in files:
			fileList.append(os.path.join(root,file))

	return fileList

def generateThumbnails(path, savePath, destPath):
	call(["php", "bin/resize.php", path, savePath, destPath])

def generate(arbo, path, thumbFolder):
	for fullPath in getFiles(path):
		print fullPath
		structName = fullPath[len(path):]
		struct = structName.split('/')
		fileName = struct[-1]
		place = None
		try:
			travel = struct[-3]
			place = struct[-2]
			placeId = ''.join(place.split(' '))
		except IndexError:
			travel = struct[-2]

		travelId = ''.join(travel.split(' '))

		if travelId not in arbo.keys():
			arbo[travelId] = {
				"title": travel,
				"places": [],
				"pics": []
			}

		tumb = False
		if place is None:
			pic = travel + '/' + fileName
			if pic not in arbo[travelId]["pics"]:
				arbo[travelId]["pics"].append(pic)
				tumb = True
		else:
			if place not in arbo[travelId]["places"]:
				arbo[travelId]["places"].append(place)

			placeIndex = arbo[travelId]["places"].index(place)
			if len(arbo[travelId]["pics"]) < placeIndex + 1:
				arbo[travelId]["pics"].append([])

			pic = travel + '/' + place + '/' + fileName
			if pic not in arbo[travelId]["pics"][placeIndex]:
				arbo[travelId]["pics"][placeIndex].append(pic)
				tumb = True

		if tumb:
			generateThumbnails(fullPath, structName, thumbFolder)

	return arbo

def main(argv):
	path = argv[0]
	jsonFile = argv[1]
	thumbFolder = argv[2]

	try:
		json_data=open(jsonFile)
		data = json.load(json_data)
		json_data.close()
	except IOError:
		data = {}
	arbo = generate(data, path, thumbFolder)
	with open(jsonFile, 'w') as outfile:
		json.dump(arbo, outfile)

if __name__ == "__main__":
	main(sys.argv[1:])
