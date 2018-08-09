var express = require('express');
var app = express();
var sql = require("mssql");

var userName = process.argv[2];
var password = process.argv[3];

// config for your database
var config = {
    user: 'sa',
    password: 'SuperSecretPassword1!',
    server: 'localhost',
    database: 'master'

};

// connect to your database
sql.connect(config, function (err) {

    if (err)  {
        console.log("Cannot connect to database? Check you credentials")

        console.log(err);
    }

    try {
        // create Request object
        var request = new sql.Request();

        // query to the database and get the records
        var selectString = 'select * from Persons where username = \''+userName+'\' AND password=\''+password+'\'';
        request.query(selectString, function (err, recordset) {

            if (err) {
                // console.log(err.message)
                if (err.message.match(/Invalid object name/i)) { // If Table does not exist, create it
                    request.query('CREATE TABLE Persons ( guid int, username varchar(255), password varchar(255))',
                        function (err, recordset) {
                            // console.log("Created Table")
                            request.query('INSERT INTO Persons (guid, username, password) VALUES (1, \'dbdbdb\', \'dbdbdb\')',
                                function (err, recordset) {
                                    // console.log("Added user dbdbdb with password dbdbdb")
                                    console.log('{"Active":"false"}')

                                })
                            setTimeout(function () {
                                sql.close()
                            }, 2000);
                    })

                }
            }
            else {
                // send records as a response
                // res.send(recordset);

                // console.log(typeof recordset.recordsets[0][0])

                if (  typeof recordset.recordsets[0][0] === "object" ) {
                    console.log('{"Active":"true"}')
                } else {
                    console.log('{"Active":"false"}')

                }


                sql.close()
            }
        })

    } catch (err) {
        // if (err instanceof Errors.NotFound)
        //     return res.status(HttpStatus.NOT_FOUND).send({message: err.message}); // 404
        // console.log(err);
        // return res.status(HttpStatus.INTERNAL_SERVER_ERROR).send({error: err, message: err.message}); // 500
    }


});

/*
CREATE TABLE Persons ( guid int, username varchar(255), password varchar(255) );

INSERT INTO Persons (guid, username, password) VALUES (1, 'dbdbdb', 'dbdbdb');

select * from Persons where username = 'dbdbdb' AND password='dbdbdb';
*/