
# Go Config Server

A simple config server written in Go, similarity to spring config server but support jsonnet as config file.

## Installation

Install by running:

```bash
go install github.com/gota33/go-config-server
```

Start server by running:
```bash
go-config-server web -http=:8080 -repo=<git repo> -user=<username> -pass=<password>
```
## Usage/Examples

Use [google/jsonnet](https://github.com/google/jsonnet) repo as example:

Start server:

```bash
go-config-server web -http=:8080 -repo=https://github.com/google/jsonnet.git
```

Query config:
```bash
curl http://localhost:8080/master/examples/arith.jsonnet
```

Response:
```json
{
   "concat_array": [
      1,
      2,
      3,
      4
   ],
   "concat_string": "1234",
   "equality1": false,
   "equality2": true,
   "ex1": 1.6666666666666665,
   "ex2": 3,
   "ex3": 1.6666666666666665,
   "ex4": true,
   "obj": {
      "a": 1,
      "b": 3,
      "c": 4
   },
   "obj_member": true,
   "str1": "The value of self.ex2 is 3.",
   "str2": "The value of self.ex2 is 3.",
   "str3": "ex1=1.67, ex2=3.00",
   "str4": "ex1=1.67, ex2=3.00",
   "str5": "ex1=1.67\nex2=3.00\n"
}
```