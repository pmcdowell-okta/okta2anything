

`docker run -e 'ACCEPT_EULA=Y' -e 'SA_PASSWORD=SuperSecretPassword1!' -p 1433:1433 -d microsoft/mssql-server-linux:2017-CU8`

`docker exec -it DOCKERPS_ID /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P 'SuperSecretPassword1!'`



