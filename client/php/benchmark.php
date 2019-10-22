<body>
<?php
require "geoipns.php";

$geodns_server = "10.105.177.129";

function getRandomIP() {
	return "".mt_rand(0,255).".".mt_rand(0,255).".".mt_rand(0,255).".".mt_rand(0,255);
}

$TEST_COUNT = 500;
if (is_null(getGeoIPNSdata("1.0.0.0", $geodns_server))) {
	$TEST_COUNT = 0;
	echo "<h2> GeoIPNS Server unreachable </h2>";
	exit;
}

// Test GeoDNS response time
$t1 = microtime(true);
$arr = array();
for ($k = 0 ; $k < $TEST_COUNT; $k++) {
	$ip = getRandomIP();
	array_push($arr, array($ip, getGeoIPNSdata($ip, $geodns_server)));
}

$t = microtime(true) - $t1;
echo "<h2>Average " . round(($t / $TEST_COUNT) * 1000000) . "&#956;s per GeoIPNS request (phpdns)" . "</h2>";

// Test with system DNS service
$t1 = microtime(true);
$arr1 = array();
for ($k = 0 ; $k < $TEST_COUNT; $k++) {
	$ip = getRandomIP();
	array_push($arr1, array($ip, getGeoIPNSdata($ip, NULL)));
}
$t = microtime(true) - $t1;
echo "<h2>Average " . round(($t / $TEST_COUNT) * 1000000) . "&#956;s per GeoIPNS request (native)" . "</h2>";

foreach (array_merge($arr, $arr1) as &$pp) {
	$ip = $pp[0]; $loc = $pp[1];
        if (is_null($loc)) {
		echo "$ip failed <br>";
	} else {
		echo "$ip ";
		print_r($loc);
		echo "<br>";
	}
}

?>
</body>

