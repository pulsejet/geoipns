<?php
require "phpdns/dns.inc.php";

function makeGeoRequest($ip, $server, $port) {
	// Answer array
	$answers = array();

	// Use PHP DNS
	if (!is_null($server)) {
		$dns_query = new DNSQuery($server, $port, 1);
		$result = $dns_query->Query($ip, "TXT");

		if (($result === false) || ($dns_query->error != 0)) {
			return null;
		}

		foreach ($result->results as $res) {
			array_push($answers, $res->data);
		}
		return $answers;
	}

	// Use system DNS service
	foreach (dns_get_record($ip, DNS_TXT) as $res) {
		array_push($answers, $res['txt']);
	}
	return $answers;
}

// Get IP information
function getGeoIPNSdata($ip, $server=NULL, $port=53) {
	// Make request
	$ip = str_replace(':', 'x', $ip) . ".geoipns.iitb.ac.in";
        $response = makeGeoRequest($ip, $server, $port);
	if (is_null($response)) return null;

	// Store KV pairs
	$answers = array();

	// Split
	foreach ($response as $res) {
		$ans = explode('=', $res);
		$answers[$ans[0]] = $ans[1];
	}

	return $answers;
}

?>
