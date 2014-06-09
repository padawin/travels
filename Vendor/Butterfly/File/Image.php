<?php

class Butterfly_File_Image extends Butterfly_File
{
	const IMAGE_PNG = 0;
	const IMAGE_JPG = 1;
	const IMAGE_GIF = 2;
	const IMAGE_JPEG_EXT = 0;
	const IMAGE_JPG_EXT = 1;
	const IMAGE_PNG_EXT = 2;
	const IMAGE_GIF_EXT = 3;

	static $allowedTypes = array('image/png', 'image/jpeg', 'image/gif');
	static $allowedExtentions = array('jpeg', 'jpg', 'png', 'gif');

	public function rotate($angle, $newImagePath = '')
	{
		if ($newImagePath == '') {
			$newImagePath = $this->_file;
		}

		$saveFunction = '';
		switch ($this->getType()) {
			case self::$allowedExtentions[Butterfly_File_Image::IMAGE_JPEG_EXT]:
			case self::$allowedExtentions[Butterfly_File_Image::IMAGE_JPG_EXT]:
			case self::$allowedTypes[Butterfly_File_Image::IMAGE_JPG]:
				$srcImage = imagecreatefromjpeg($this->_file);
				$saveFunction = 'imagejpeg';
				break;
			case self::$allowedExtentions[Butterfly_File_Image::IMAGE_PNG_EXT]:
			case self::$allowedTypes[Butterfly_File_Image::IMAGE_PNG]:
				$srcImage = imagecreatefrompng($this->_file);
				$saveFunction = 'imagepng';
				break;
			case self::$allowedExtentions[Butterfly_File_Image::IMAGE_GIF_EXT]:
			case self::$allowedTypes[Butterfly_File_Image::IMAGE_GIF]:
				$srcImage = imagecreatefromgif($this->_file);
				$saveFunction = 'imagegif';
				break;
			default:
				throw new Exception('Unknown image type');
				break;
		}

		$newImage = imagerotate($srcImage, $angle, 0);

		$saveFunction($newImage, $newImagePath);

		//delete tmp images
		imagedestroy($newImage);
		imagedestroy($srcImage);

		$this->_file = $newImagePath;
	}

	/**
	 *
	 * Resize the picture with the given options. The ratio between
	 * width and height is not changed, whatever is the option
	 *
	 * @param $args options to define how to resize:
	 *	  percent:
	 *	  width:
	 *	  height:
	 *	  crop:
	 * @param $newImagePath if given, the resized picture will be saved
	 * there, instead of modificate the original picture
	 *
	 */
	public function resize($args = array(), $newImagePath = '')
	{
		if ($newImagePath == '') {
			$newImagePath = $this->_file;
		}

		//check args
		if (
			isset($args['crop']) && $args['crop'] == null &&
			isset($args['width']) && $args['width'] == null &&
			isset($args['height']) && $args['height'] == null &&
			isset($args['percent']) && $args['percent'] == null
		) {
			throw new Exception(
				'You must at least provide the new height (heightMax),
				the new width (widthMax) or the resize percent (percent)
				of the image in the first arg of the function'
			);
		}

		if (!isset($args['enlarged']) || $args['enlarged'] == true) {
			$canBeEnlarged = true;
		}
		else {
			$canBeEnlarged = false;
		}

		if (!isset($args['crop']) || $args['crop'] == false) {
			$crop = false;
		}
		else {
			$crop = true;
		}

		if (isset($args['percent']) && $args['percent'] != null && $args['percent'] <= 0) {
			throw new Exception('The percent must be higher than 0');
		}

		if (
			isset($args['width']) && $args['width'] != null && !(is_numeric($args['width']) && (int) $args['width'] = $args['width'])
			|| isset($args['height']) && $args['height'] != null && !(is_numeric($args['height']) && (int) $args['heigthMax'] = $args['height'])
		) {
			throw new Exception('The max width and height must be integers');
		}

		$saveFunction = '';
		switch ($this->getType()) {
			case self::$allowedExtentions[Butterfly_File_Image::IMAGE_JPEG_EXT]:
			case self::$allowedExtentions[Butterfly_File_Image::IMAGE_JPG_EXT]:
			case self::$allowedTypes[Butterfly_File_Image::IMAGE_JPG]:
				$srcImage = imagecreatefromjpeg($this->_file);
				$saveFunction = 'imagejpeg';
				break;
			case self::$allowedExtentions[Butterfly_File_Image::IMAGE_PNG_EXT]:
			case self::$allowedTypes[Butterfly_File_Image::IMAGE_PNG]:
				$srcImage = imagecreatefrompng($this->_file);
				$saveFunction = 'imagepng';
				break;
			case self::$allowedExtentions[Butterfly_File_Image::IMAGE_GIF_EXT]:
			case self::$allowedTypes[Butterfly_File_Image::IMAGE_GIF]:
				$srcImage = imagecreatefromgif($this->_file);
				$saveFunction = 'imagegif';
				break;
			default:
				throw new Exception('Unknown image type');
				break;
		}

		//get the new size
		$oldWidth = imagesx($srcImage);
		$oldHeight = imagesy($srcImage);

		//define new size from old ones and from the args
		if (isset($args['percent']) && $args['percent'] != null) {
			$size = $this->_getSizeFromPercent($oldWidth, $oldHeight, $args['percent']);
		}
		elseif (isset($args['width']) || isset($args['height'])) {
			if (isset($args['width'])) {
				$width = $args['width'];
			}
			else {
				$width = null;
			}

			if (isset($args['height'])) {
				$height = $args['height'];
			}
			else {
				$height = null;
			}
			$size = $this->_getSizeFromMaxSize($oldWidth, $oldHeight, $width, $height, $canBeEnlarged, $crop);
		}
		else {
			$size = array(
				'width' => $oldWidth,
				'height' => $oldHeight
			);
		}
		$newWidth = $size['width'];
		$newHeight = $size['height'];

		//define the coordinate to create the new picture, if a crop is asked
		$srcX = 0;
		$srcY = 0;

		//a crop can only be done if the 2 output size are given
		if ($crop && !empty($width) && !empty($height)) {
			if ($newWidth > $width) {
				$croppedWidth = $oldWidth * $width / $newWidth;
				$srcX = ($oldWidth - $croppedWidth) / 2;
				$oldWidth = $croppedWidth;
			}
			elseif ($newHeight > $height) {
				$croppedHeight = $oldHeight * $height / $newHeight;
				$srcY = ($oldHeight - $croppedHeight) / 2;
				$oldHeight = $croppedHeight;
			}

			$newWidth = $width;
			$newHeight = $height;
		}

		$newImage = imagecreatetruecolor($newWidth, $newHeight);
		imagecopyresampled($newImage, $srcImage, 0, 0, $srcX, $srcY, $newWidth, $newHeight, $oldWidth, $oldHeight);

		$saveFunction($newImage, $newImagePath);

		//delete tmp images
		imagedestroy($newImage);
		imagedestroy($srcImage);

		$this->_file = $newImagePath;
	}


	/**
	 *
	 * @return array containing the new width and height
	 *
	 */
	private function _getSizeFromPercent($width, $height, $percent)
	{
		$width = $percent * $width / 100;
		$height = $percent * $height / 100;

		return array(
			'height' => $height,
			'width' => $width
		);
	}

	/**
	 *
	 * calculate the new width and height from the widthMax, the heightMax
	 * or both of them
	 *
	 * @return array containing the new width and height
	 *
	 */
	private function _getSizeFromMaxSize($widthImage, $heightImage, $widthMax = null, $heightMax = null, $canBeEnlarged = true, $crop = false)
	{
		$size = array();
		$coefWidth = 1;
		$coefHeight = 1;

		if ($widthMax == null && $heightMax == null) {
			throw new Exception('At least the max height or the max width must be provided');
		}
		//only height given
		elseif ($widthMax == null) {
			$widthMax = 0;
			//coef to invert the ratios, to avoid to have a 0*0 picture
			$coefHeight = -1;
		}
		//only width given
		elseif ($heightMax == null) {
			$heightMax = 0;
			//coef to invert the ratios, to avoid to have a 0*0 picture
			$coefWidth = -1;
		}

		$ratioWidth = $widthMax / $widthImage * $coefWidth;
		$ratioHeight = $heightMax / $heightImage * $coefHeight;

		if (!$canBeEnlarged && $ratioWidth >= 1 && $ratioHeight >= 1) {
			$size['width'] = $widthImage;
			$size['height'] = $heightImage;
		}
		elseif (!$crop && $ratioHeight <= $ratioWidth || $crop && $ratioHeight > $ratioWidth) {
			$size['height'] = $heightMax;
			$size['width'] = $widthImage * $heightMax / $heightImage;
		}
		else {
			$size['width'] = $widthMax;
			$size['height'] = $heightImage * $widthMax / $widthImage;
		}

		return $size;
	}

	public function hasCorrectType()
	{
		return in_array($this->getType(), self::$allowedTypes);
	}

	/**
	 *
	 * echo content of the file
	 *
	 */
	public function __toString()
	{
		echo base64_encode(file_get_contents($this->_file));
	}

	public static function createFromBase64($base64File, $pathToSave)
	{
		$file = explode(',', $base64File);
		preg_match('/image\/([\w]+)/', $file[0], $match);

		//check type
		if (
			!isset($match[1]) ||
			(
				!in_array($match[1], self::$allowedTypes) &&
				!in_array($match[1], self::$allowedExtentions)
			)
		) {
			throw new Exception('Unknown Type');
		}

		$image = imagecreatefromstring(base64_decode(str_replace(' ', '+', $file[1])));
		$function = 'image' . $match[1];
		$function($image, $pathToSave);

		return new self($pathToSave);
	}
}
