/////////////////////////////////////////////////////////////
// Setup:
//      On Linux,
//          service start mongodb
//      Install NPM:
//          npm install mongodb@2.2.5 --save
//      Then go to Mongo Console:
//          use users;
//          db.users.insert({"username":"usertest","password":"usertest"});
//
// Test:
//  node index.js usertest usertest

const MongoClient = require('mongodb').MongoClient;
const url = "mongodb://localhost:27017/users" // Name of database is users

var args = process.argv.slice(2);

var username = args[0]
var password = args[1]

var found = false;

MongoClient.connect(url, {useNewUrlParser: true}, function (err, db) {
    if (err) throw err;
    let dbo = db.db("users"); //name of database.. might be redundant, dunno

    var stream = dbo.collection("users").find({"username": username, "password": password}).stream() //name of collection is users
    stream.on('error', function (err) {
        // console.error(err)
        console.log('{"Active":"false"}')
        db.close();
    })
    stream.on('data', function (doc) {
        // console.log(doc)
        found = true;
        db.close()
    })
        .on('end', function () {
            // final callback
            if (found == true) {
                console.log('{"Active":"true"}')

            } else {
                console.log('{"Active":"false"}')

            }
            db.close();
        });
});



