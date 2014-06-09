#!/usr/bin/env python

import sys
import os
import json

def getFiles(path):
	fileList = []
	rootdir = sys.argv[1]
	for root, subFolders, files in os.walk(path):
		for file in files:
			fileList.append(os.path.join(root,file))

	return fileList

def generate(path):
	arbo = {}
	for f in getFiles(path):
		f = f[len(path):]
		struct = f.split('/')
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
			arbo[travelId]["pics"].append(fileName)
		else:
			if place not in arbo[travelId]["places"]:
				arbo[travelId]["places"].append(place)

			placeIndex = arbo[travelId]["places"].index(place)
			if len(arbo[travelId]["pics"]) < placeIndex + 1:
				arbo[travelId]["pics"].append([])

			arbo[travelId]["pics"][placeIndex].append(fileName)
	return arbo

def main(argv):
	path = argv[0]
	destFile = argv[1]
	arbo = generate(path)
	with open(destFile, 'w') as outfile:
		json.dump(arbo, outfile)

if __name__ == "__main__":
	main(sys.argv[1:])
