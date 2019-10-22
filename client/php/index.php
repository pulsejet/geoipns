<!DOCTYPE HTML>
<html>
<head>
	<title> GeoIPNS Lookup </title>
</head>

<body>
	<h1> GeoIPNS Lookup </h1>

	<form action="index.php" method="post">
		IP: <input type="text" name="ip" placeholder="Enter IP Address"><br>
		<input type="submit">
	</form>

<?php
require "geoipns.php";

$geodns_server = "10.105.177.129";

// Get POST data
if ($_SERVER['REQUEST_METHOD'] === 'POST') {
	$ip = $_POST["ip"];
	$loc = getGeoIPNSdata($ip, $geodns_server);
        if (is_null($loc)) {
		echo "Query - $ip failed <br>";
	} else {
		print_r($loc);
	}
}

?>
</body>
</html>

