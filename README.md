ElSimulator
========

(./elSimulator -help)
```
Usage of ./elSimulator:
  -baseDirectory="elSimulatorCurrent": directory with file to read (elSimulatorCurrent to use directory elSimulator)
  -bindingAddress="localhost:4000": The binding address
  -parameterRegex=".*": Parameter regex
  -proxyAddress="http://localhost:4000/file": The binding address
```

##Start server ./elSimulator

####Feature Reader

|URL|file|http code|
|---|---|---|
|http://localhost:4000/file/test|current base/file/test/GET/withoutParameter.xml|500|
|http://localhost:4000/file/test?param=value|current base/file/test/GET/param_value.json|203|


###Feature proxy (proxyAddress=http://www.google.fr)

|URL|called url|file|
|---|---|---|
|http://localhost:4000/proxy/?q=test|http://www.google.fr/?q=test|current base/proxy/GET/q_test|