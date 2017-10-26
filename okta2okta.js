var http = require("https");

var args = process.argv.slice(2);

var username = args[0]
var password = args[1]

var options = {
  "method": "POST",
  "hostname": "companyx.okta.com",
  "port": null,
  "path": "/api/v1/authn",
  "headers": {
    "accept": "application/json",
    "content-type": "application/json",
    "cache-control": "no-cache",
    "postman-token": "d3fefc1d-89e0-a411-02e2-c4d30c9abda0"
  }
};

var req = http.request(options, function (res) {
  var chunks = [];

  res.on("data", function (chunk) {
    chunks.push(chunk);
  });

  res.on("end", function () {
    var body = Buffer.concat(chunks);
    body=body.toString()

      if ( body.indexOf("Authentication failed") >0 ) {

          console.log('{"Active":"false"}')
      } else {
          console.log('{"Active":"true"}')
      }


      // console.log(body);
  });
});

req.write(JSON.stringify({ username: username,
  password: password,
  options: 
   { multiOptionalFactorEnroll: true,
     warnBeforePasswordExpired: true } }));
req.end();

