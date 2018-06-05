#!/usr/bin/env python3

import sys
import os
import json
import errno
import time
from images import process_image


def get_files(path):
    file_list = []
    for root, sub_folders, files in os.walk(path):
        for file_name in sorted(files):
            file_list.append(os.path.join(root, file_name))

    return file_list


def create_dir(file_path):
    dirname = os.path.dirname(file_path)
    if os.path.exists(dirname):
        return

    try:
        os.makedirs(dirname)
    except OSError as exc:  # Guard against race condition
        if exc.errno != errno.EEXIST:
            raise


def generate_thumbnails(path, save_path, dest_path):
    for i, im in enumerate(process_image(path), start=1):
        destination = "{directory}/{width}x{height}x{crop}/{path}".format(
            directory=dest_path,
            width=im[0]['width'],
            height=im[0]['height'],
            crop=int(im[0]['crop']),
            path=save_path
        )
        create_dir(destination)
        im[1].save(destination, "JPEG", quality=95, optimize=True)


def generate(arbo, path, thumb_folder):
    for full_path in get_files(path):
        print(full_path)
        struct_name = full_path[len(path):]
        struct = struct_name.split('/')
        file_name = struct[-1]
        place = None
        try:
            travel = struct[-3]
            place = struct[-2]
        except IndexError:
            travel = struct[-2]

        travel_id = ''.join(travel.split(' '))

        if travel_id not in arbo.keys():
            arbo[travel_id] = {
                "title": travel,
                "places": [],
                "pics": []
            }

        tumb = False
        if place is None:
            pic = travel + '/' + file_name
            if pic not in arbo[travel_id]["pics"]:
                arbo[travel_id]["pics"].append(pic)
                tumb = True
        else:
            if place not in arbo[travel_id]["places"]:
                arbo[travel_id]["places"].append(place)

            place_index = arbo[travel_id]["places"].index(place)
            if len(arbo[travel_id]["pics"]) < place_index + 1:
                arbo[travel_id]["pics"].append([])

            pic = travel + '/' + place + '/' + file_name
            if pic not in arbo[travel_id]["pics"][place_index]:
                arbo[travel_id]["pics"][place_index].append(pic)
                tumb = True

        if tumb:
            generate_thumbnails(full_path, struct_name, thumb_folder)

    return arbo


def main(argv):
    path = argv[0]
    json_file = argv[1]
    thumb_folder = argv[2]

    try:
        json_data = open(json_file)
        data = json.load(json_data)
        json_data.close()
    except IOError:
        data = {}
    arbo = generate(data, path, thumb_folder)
    with open(json_file, 'w') as outfile:
        json.dump(arbo, outfile)


if __name__ == "__main__":
    main(sys.argv[1:])
