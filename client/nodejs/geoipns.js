/**
=================
Sample usage:
=================

const getGeoIPNS = require('./geoipns.js');

const config = {
    servers: ['10.105.177.129'],
    suffix: "geoipns.iitb.ac.in",
};

getGeoIPNS('1.1.1.1', config).then(res => {
    console.log(res)
}).catch(err => {
    console.error(err)
});

*/

const dns = require('dns');

/** Request for location and ASN data with no time limit */
function getGeoIPNS(ip, config) {
    // Set servers
    const resolver = new dns.Resolver();
    if (config.servers) {
        resolver.setServers(config.servers);
    }

    return new Promise((resolve, reject) => {
        const result = {};
        resolver.resolveTxt(`${ip}.${config.suffix || "geoipns"}`, (err, res) => {
            if (err) {
                reject(err);
                return;
            }
            res.forEach(r => {
                const record = r[0].split('=');
                result[record[0]] = record[1]
            });
            resolve(result);
        });

        // Set timeout
        setTimeout(() => reject(null), config.timeout || 1000);
    });
}

module.exports = getGeoIPNS;
