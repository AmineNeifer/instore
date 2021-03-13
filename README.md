## ``Instore`` Service and Client 

- ``Instore`` is a microservice which is responsible for managing an in-memory
key-value store. The user of that service can add, retrieve, remove... key-value pairs via that service. 

- ``Instore`` service and the consumer of that service are implemented in both ``Go`` language.

## Service Definition 

```proto
package store;

option go_package="storepb";

message AddCsvRequest {
    string key = 1;
    string value = 2;
}

message AddCsvResponse {
    string result = 1;
}

message RemoveCsvRequest {
    string key = 1;
    string value = 2;
}

message RemoveCsvResponse {
    string result = 1;
}

message GetvCsvRequest {
    string key = 1;
}

message GetvCsvResponse {
    string result = 1;
}

message GetkCsvRequest {
    string value = 2;
}

message GetkCsvResponse {
    string result = 1;
}

message RemoveAllCsvRequest {
    string msg = 1;
}

message RemoveAllCsvResponse {
    string result = 1;
}

message AddCsvFromFileRequest {
    string key = 1;
    string value = 2;
}
message AddCsvFromFileResponse {
    string result = 1;
}

message RemoveCsvFromFileRequest {
    string key = 1;
    string value = 2;
}

message RemoveCsvFromFileResponse {
    string result = 1;
}

message GetAllCsvRequest {
    string msg = 1;
}

message GetAllCsvResponse {
    string key = 1;
    string value = 2;
}

message UseCsvRequest {
    string msg = 1;
}

message UseCsvResponse {
    string result = 1;
}

service StoreService {
    // Unary
    rpc AddCsv(AddCsvRequest) returns (AddCsvResponse) {};
    rpc RemoveCsv(RemoveCsvRequest) returns (RemoveCsvResponse) {};
    rpc GetvCsv(GetvCsvRequest) returns (GetvCsvResponse) {};
    rpc GetkCsv(GetkCsvRequest) returns (GetkCsvResponse) {};
    rpc RemoveAllCsv(RemoveAllCsvRequest) returns (RemoveAllCsvResponse) {};
    rpc UseCsv(UseCsvRequest) returns (UseCsvResponse) {};
    // client streaming
    rpc AddCsvFromFile(stream AddCsvFromFileRequest) returns (AddCsvFromFileResponse) {};
    rpc RemoveCsvFromFile(stream RemoveCsvFromFileRequest) returns (RemoveCsvFromFileResponse) {};
    // server streaming
    rpc GetAllCsv(GetAllCsvRequest) returns (stream GetAllCsvResponse) {};
}


message Data {
    string id = 1;
    string key = 2;
    string value = 3;
}
message AddDbRequest {
    Data data = 1;
}

message AddDbResponse {
    Data data = 1;
}

message UseDbRequest {
    string msg = 1;
}
message UseDbResponse {
    string result = 1;
}

message RemoveDbRequest {
    string key = 1;
    string value = 2;
}

message RemoveDbResponse {
    string result = 1;
}

message GetvDbRequest {
    string key = 1;
}
message GetvDbResponse {
    Data data = 1;
}

service StoreDbService {
    rpc AddDb(AddDbRequest) returns (AddDbResponse) {};
    rpc GetvDb(GetvDbRequest) returns (GetvDbResponse) {};
    rpc RemoveDb(RemoveDbRequest) returns (RemoveDbResponse) {};
    rpc UseDb(UseDbRequest) returns (UseDbResponse) {};
}
```
## How to use
### Benefit of the mongoDB feature
In order to start mongoDB and benefit from its features run the following
command from the root directory (instore/)
```
mongodb --dbpath storage/db
```

### Building and Running Server

In order to build, Go to location (instore/server) and execute the following
 shell command,
```
go build -o bin/server
```

In order to run, Go to location (instore/server) and execute the following
shell command,

```
./bin/server
```

You can also run it from the root directory location (instore/)
```
go run server/server.go
```

### Building and Running Client
In order to build, Go to location (instore/command) and execute the following
 shell command,
```
go build -i -v -o bin/command
```

In order to run, Go to location (instore/command) and execute the following
shell command,

```
./bin/command
```

You can also run it from the root directory location(instore/)
```
go run command/command.go
```

You will be prompted to a ``CLI``, that is able to invoke function from 
``instore/client/client.go`` here is an example of the experience
```
[cmd]  Hello! I am excited to have you as a user :D Enjoy!
[cmd]  You could use any of these commands
[cmd]  ----------------------------------------------------------------------------------------------
[cmd]  add       <key>  <value>  ------- to add a key-value pair to the store
[cmd]  get[v]    <key>           ------- to get existing values corresponding to a key from the store
[cmd]  getk      <value>         ------- to get existing keys corresponding to a value from the store
[cmd]  remove    <key>  <value>  ------- to remove a key-value pair from the store
[cmd]  exit/quit                 ------- to quit the program
[cmd]  ----------------------------------------------------------------------------------------------
[usr]> add tunis france
[srv]  Added: {'tunis': 'france'}
[usr]> add amine zenly
[srv]  Added: {'amine': 'zenly'}
[usr]> add holberton school
[srv]  Added: {'holberton': 'school'}
[usr]> use db
[srv]  Using MongoDB now...
[usr]> add betty holberton
data:<id:"604bf6493d024033948a623f" key:"betty" value:"holberton" > 
[usr]> exit
[cmd]  Bye Bye! :)
```


``[cmd] `` is output generated by the command.go file  
``[usr]>`` is prefix to show that it is waiting for the user input  
``[svr] `` is output generated by the server.go file  
``[clt] `` is output generated by the client.go  

## Directories
`client` contains client.go which is client-side code (only invoked from CLI)  
`command` contains command.go which is command line interface code  
`server` contains server.go which is server-side code  
`storage` contains storage.go which has functions to deal with csv files  
`storepb` contains store.proto and store.pb.go  

## Additional Information
### Commands
#### ``use``
``[usr]> use db`` enables MOngoDB mode (stores files and retreives files from mongodb) 
#### ``add``
``[usr]> add Tunisia tunis`` adds ``Tunsia-tunis`` key-value pair to ``data.csv`` 
#### ``addcsv``
``[usr]> add new.csv`` adds ``new.csv`` key-value pairs to ``data.csv`` 
#### ``get / getv``
``[usr]> get Tunisia`` prints values paired to the key ``Tunsia`` in ``data.csv`` 
#### ``getall``
``[usr]> getall`` prints all key-value pairs in ``data.csv`` 
#### ``getk``
``[usr]> getk tunis`` prints keys paired with the value ``tunis`` in ``data.csv`` 
#### ``remove / delete``
``[usr]> remove Tunisia tunis`` removes ``Tunsia-tunis`` key-value pair from ``data.csv`` 
#### ``removecsv``
``[usr]> removecsv new.csv`` remove key-value pairs from ``data.csv`` that are present in ``new.csv`` 
### Csv files used in the service
``[usr]> data.csv`` is generated in the same directory from which you executed server.go
if you generated a new csv file from which you want to add to data.csv with ``addcsv`` command
the file needs to be in the directory in which you ran command.go 

### Generate Server and Client side code 
Pre-generated stub file is included in the go project. If you need to generate the stub files please use the below
 command from the root directory(inside instore directory)
``` 
protoc storepb/store.proto --go_out=plugins=grpc:.
``` 