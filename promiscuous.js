//Script always returns true, so anyone can login
var args = process.argv.slice(2);

var username = args[0]
var password = args[1]

console.log('{"Active":"true"}')
