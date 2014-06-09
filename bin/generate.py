#!/usr/bin/env python

import sys
import os
from os.path import isfile, join

def getFiles(path):
	fileList = []
	rootdir = sys.argv[1]
	for root, subFolders, files in os.walk(path):
		for file in files:
			fileList.append(os.path.join(root,file))

	return fileList

def main(argv):
	path = argv[0]

	fileList = getFiles(path)
	print fileList

if __name__ == "__main__":
	main(sys.argv[1:])
