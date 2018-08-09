

`docker run -e 'ACCEPT_EULA=Y' -e 'SA_PASSWORD=SuperSecretPassword1!' -p 1433:1433 -d microsoft/mssql-server-linux:2017-CU8`

`docker exec -it DOCKERPS_ID /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P 'SuperSecretPassword1!'`

``` 
CREATE TABLE Persons ( guid int, username varchar(255), password varchar(255) );

INSERT INTO Persons (guid, username, password) VALUES (1, 'dbdbdb', 'dbdbdb');

select * from Persons where username = 'dbdbdb' AND password='dbdbdb';
```


