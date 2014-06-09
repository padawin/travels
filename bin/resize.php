<?php

require("Vendor/Butterfly/File.php");
require("Vendor/Butterfly/File/Image.php");


$thumbSizes = array(
	array('width' => 100, 'height' => 100, 'crop' => 1),
	array('width' => 118, 'height' => 133, 'crop' => 1),
	array('width' => 118, 'height' => 133, 'crop' => 0),
	array('width' => 250, 'height' => 250, 'crop' => 1),
	array('width' => 1024, 'height' => 768, 'crop' => 0)
);

$imageLocation = $argv[1];
$imageSavePath = $argv[2];
$destinationPath = $argv[3];

//copy the picture to the htdocs folder
foreach ($thumbSizes as $size) {
	$image = new Butterfly_File_Image($imageLocation);
	$dirname = dirname($destinationPath . '/' . implode('x', $size) . '/' . $imageSavePath);
	if (!is_dir($dirname)) {
		mkdir($dirname, 0777, true);
	}

	$file = $image->getFileName();
	$exifs = exif_read_data($file);

	if (isset($exifs['Orientation'])) {
		if ($exifs['Orientation'] == 8) {
			$image->rotate(90, '/tmp/img');
		}
		elseif ($exifs['Orientation'] == 6) {
			$image->rotate(270, '/tmp/img');
		}
	}

	$image->resize($size, uniqid('/tmp/img'));
	$image->copy($destinationPath . '/' . implode('x', $size) . '/' . $imageSavePath);
}
