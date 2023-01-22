<?php

declare(strict_types=1);

function listDir(string $path) : array{
	$content = array_diff(scandir($path), array('.', '..'));
	$res = [];
	foreach ($content as $item) {
		$fullPath = $path . "/" . $item;
		if (is_dir($fullPath)) {
			$res = array_merge($res, listDir($fullPath));
		}
		else if ($item != "random.php") {
			array_push($res, $fullPath);
		}
	}

	return $res;
}

$path = realpath(__DIR__);
$destWidth = {{ .DestWidth }};
$destHeight = {{ .DestHeight }};

// load n random images from current directory
$files = listDir($path);
shuffle($files);
$files = array_slice($files, 0, {{ .CountThumbs }});

// Debug
//fwrite(STDERR, $files[0]."\n");
//fwrite(STDERR, $files[1]."\n");
//fwrite(STDERR, $files[2]."\n");
//fwrite(STDERR, $files[3]."\n");

// Combine images in a grid:
$dstImage = imagecreatetruecolor($destWidth, $destHeight);
{{ range $i, $thumb := .ThumbPositions }}
	$srcImage = imagecreatefromjpeg($files[{{ $i }}]);
	imagecopy($dstImage, $srcImage, {{ $thumb.X }}, {{ $thumb.Y }}, 0, 0, {{ $.ThumbWidth }}, {{ $.ThumbHeight }});
{{ end }}

// display image
header('Content-type: image/jpeg');
imagejpeg($dstImage);
