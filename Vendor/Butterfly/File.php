<?php

class Butterfly_File
{

	/**
	 *
	 * string, absolute file path
	 *
	 */
	protected $_file;

	/**
	 *
	 * Construct
	 * Set the filepath
	 *
	 */
	public function __construct($filePath = '')
	{
		if ($filePath != '' && ! is_file($filePath)) {
			throw new Exception('The given file does not exist');
		}

		$this->_file = $filePath;
	}

	/**
	 *
	 * return the name of the file
	 *
	 */
	public function getFileName()
	{
		return $this->_file;
	}

	/**
	 *
	 * Return the mime type of the file
	 *
	 * @return string
	 *
	 */
	public function getType()
	{
		if (function_exists('mime_content_type')) {
			return mime_content_type($this->_file);
		}
		elseif (function_exists('finfo_open')) {
			$finfo = finfo_open(FILEINFO_MIME_TYPE);
			return finfo_file($finfo, $this->_file);
		}
		else {
			$type = explode('.', $this->_file);
			return $type[count($type) - 1];
		}
	}

	/**
	 *
	 * return the weight of the file
	 *
	 * @return int
	 *
	 */
	public function getWeight()
	{
		return filesize($this->_file);
	}

	/**
	 *
	 * Move this->_file into another file
	 *
	 */
	public function move($newName, $uploaded = false)
	{
		if ($uploaded) {
			move_uploaded_file($this->_file, $newName);
		}
		else {
			rename($this->_file, $newName);
		}

		$this->_file = $newName;
	}

	/**
	 *
	 * Copy $this->_file to $newName
	 *
	 * @return Butterfly_File the created file
	 *
	 *
	 */
	public function copy($newName)
	{
		copy($this->_file, $newName);
		$type = get_class($this);

		return new $type($newName);
	}

	/**
	 *
	 * echo content of the file
	 *
	 */
	public function __toString()
	{
		echo $this->_file;
	}

	/**
	 *
	 * Recursively delete a directory
	 *
	 */
	public static function rmDir($dir)
	{
		if (is_dir($dir)) {
			$elements = scandir($dir);
			foreach ($elements as $element) {
				if ($element != '.' && $element != '..') {
					if (is_dir($dir . '/' . $element)) {
						rmdir($dir . '/' . $element);
					}
					else {
						unlink($dir . '/' . $element);
					}
				}
			}
			reset($elements);
			rmdir($dir);
		}
	}
}
