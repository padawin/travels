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

def generate(path, thumbFolder):
	arbo = {}
	for fullPath in getFiles(path):
		structName = fullPath[len(path):]
		generateThumbnails(fullPath, structName, thumbFolder)
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

		if place is None:
			arbo[travelId]["pics"].append(travel + '/' + fileName)
		else:
			if place not in arbo[travelId]["places"]:
				arbo[travelId]["places"].append(place)

			placeIndex = arbo[travelId]["places"].index(place)
			if len(arbo[travelId]["pics"]) < placeIndex + 1:
				arbo[travelId]["pics"].append([])

			arbo[travelId]["pics"][placeIndex].append(travel + '/' + place + '/' + fileName)
	return arbo

def main(argv):
	path = argv[0]
	jsonFile = argv[1]
	thumbFolder = argv[2]
	arbo = generate(path, thumbFolder)
	with open(jsonFile, 'w') as outfile:
		json.dump(arbo, outfile)

if __name__ == "__main__":
	main(sys.argv[1:])
