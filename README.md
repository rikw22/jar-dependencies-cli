# jar-dependencies-cli
Returns a json with dependencies compiled into a jar/war file



### Usage

```bash
./jar-dependencies -f test.war | jq -c '.[] | select( .Name | contains("spring-core") )'
```

returns
```json
{
    "Name":"spring-core",
    "Version":"5.3.4",
    "FullName":"spring-core-5.3.4.jar"
}
```