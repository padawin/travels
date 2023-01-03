#!/usr/bin/env python3

import sys
import os
import json
import errno
from images import process_image_for_config


def get_files(path):
    file_list = []
    for root, sub_folders, files in os.walk(path):
        for file_name in sorted(files):
            file_list.append(os.path.join(root, file_name))

    file_list.sort()
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


def generate_thumbnail(source_dir, path, image_format, dest_dir):
    width, height, crop = image_format.split("x")
    input_file = "{}/{}".format(source_dir, path)
    im = process_image_for_config(
        input_file,
        {"width": int(width), "height": int(height), "crop": bool(int(crop))}
    )
    destination = "{directory}/{image_format}/{path}".format(
        directory=dest_dir,
        image_format=image_format,
        path=path
    )
    create_dir(destination)
    im.save(destination, "JPEG", quality=95, optimize=True)


def generate_arbo(arbo, path):
    arbo['latest'] = list()
    if 'travels' not in arbo:
        arbo['travels'] = dict()
    for full_path in get_files(path):
        struct_name = full_path[len(path):]
        struct = struct_name.strip("/").split('/')
        file_name = struct[-1]
        place = None
        try:
            travel = struct[-3]
            place = struct[-2]
        except IndexError:
            travel = struct[-2]

        travel_id = ''.join(travel.split(' '))

        if travel_id not in arbo['travels'].keys():
            arbo['travels'][travel_id] = {
                "title": travel,
                "places": [],
                "pics": []
            }

        if place is None:
            pic = travel + '/' + file_name
            if pic not in arbo['travels'][travel_id]["pics"]:
                arbo['travels'][travel_id]["pics"].append(pic)
                arbo['latest'].append(pic)
        else:
            if place not in arbo['travels'][travel_id]["places"]:
                arbo['travels'][travel_id]["places"].append(place)

            place_index = arbo['travels'][travel_id]["places"].index(place)
            if len(arbo['travels'][travel_id]["pics"]) < place_index + 1:
                arbo['travels'][travel_id]["pics"].append([])

            pic = travel + '/' + place + '/' + file_name
            if pic not in arbo['travels'][travel_id]["pics"][place_index]:
                arbo['travels'][travel_id]["pics"][place_index].append(pic)
                arbo['latest'].append(pic)

    return arbo


def _parse_source_directory(source_dir, json_file):
    try:
        json_data = open(json_file)
        data = json.load(json_data)
        json_data.close()
    except IOError:
        data = {}
    try:
        data = generate_arbo(data, source_dir)
    except KeyboardInterrupt:
        pass
    return data


if __name__ == "__main__":
    operation = sys.argv[1]
    argv = sys.argv[2:]
    if operation == "json":
        json_file = argv[0]
        source_dir = argv[1]
        data = _parse_source_directory(source_dir, json_file)
        with open(json_file, 'w') as outfile:
            json.dump(data, outfile, sort_keys=True)
    elif operation == "thumb":
        source_dir = argv[0]
        file_path = argv[1]
        image_format = argv[2]
        dest_dir = argv[3]

        generate_thumbnail(source_dir, file_path, image_format, dest_dir)
