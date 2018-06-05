import functools
from PIL import Image

images_config = [
    {"crop": False, "height": 133, "width": 118},
    {"crop": False, "height": 768, "width": 1024},
    {"crop": True, "height": 100, "width": 100},
    {"crop": True, "height": 133, "width": 118},
    {"crop": True, "height": 250, "width": 250},
]


def image_transpose_exif(im):
    """
    Applies EXIF transformations to `im` (if present) and returns
    the result as a new image.
    Taken from [here](https://stackoverflow.com/a/30462851/1667018)
    and adapted.
    """
    exif_orientation_tag = 0x0112  # contains an integer, 1 through 8
    exif_transpose_sequences = [  # corresponding to the following
        [],
        [Image.FLIP_LEFT_RIGHT],
        [Image.ROTATE_180],
        [Image.FLIP_TOP_BOTTOM],
        [Image.FLIP_LEFT_RIGHT, Image.ROTATE_90],
        [Image.ROTATE_270],
        [Image.FLIP_TOP_BOTTOM, Image.ROTATE_90],
        [Image.ROTATE_90],
    ]

    try:
        seq = exif_transpose_sequences[im._getexif()[exif_orientation_tag] - 1]
    except Exception:
        return im
    else:
        if seq:
            return functools.reduce(lambda im, op: im.transpose(op), seq, im)
        else:
            return im


def process_image_for_config(fname, config):
    def new_img(size, im_to_paste=None):
        im = Image.new("RGB", size)
        im.paste(
            im_to_paste,
            (
                (size[0] - im_to_paste.size[0]) // 2,
                (size[1] - im_to_paste.size[1]) // 2,
            ),
        )
        return im
    im = image_transpose_exif(Image.open(fname))
    size = (config["width"], config["height"])
    if config["crop"]:
        im = im.resize(size, Image.ANTIALIAS)
        im = new_img(size, im)
    else:
        scale_x = size[0] / im.size[0]
        scale_y = size[1] / im.size[1]
        scale = min(scale_x, scale_y)
        size = tuple(int(round(value * scale)) for value in im.size)
        im = new_img(size, im.resize(size, Image.ANTIALIAS))
    return im


def process_image(fname):
    for config in images_config:
        yield config, process_image_for_config(fname, config)
