//Script always returns true, so anyone can login
var args = process.argv.slice(2);

var username = args[0]
var password = args[1]



console.log('{"Active":"true","telephoneNumber":"pre'+makeid(20)+'","departmentNumber":"pre'+makeid(20)+'"}')
// console.log('{"Active":"true","guidguid":"111111"}')

function makeid(length) {
    var result           = '';
    var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    var charactersLength = characters.length;
    for ( var i = 0; i < length; i++ ) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }
    return result;
}
