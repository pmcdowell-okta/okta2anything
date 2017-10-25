var fs = require('fs'); 
var parse = require('csv-parse');

var userName = process.argv[2];
var password = process.argv[3];

var csvData=[];
fs.createReadStream("users.csv")
    .pipe(parse({delimiter: ','}))
    .on('data', function(csvrow) {
//        console.log(csvrow);
        //do something with csvrow
        csvData.push(csvrow);        
    })
    .on('end',function() {

        var DudeFoundMyCarAndUser=0;
        for (var item of csvData) {
            if (item[0]==userName && item[1]==password) {
                DudeFoundMyCarAndUser=1
            }

        }

        if (DudeFoundMyCarAndUser!=0) {
            console.log('{"Active":"true"}')
        } else {
            console.log('{"Active":"false"}')
        }

    });

