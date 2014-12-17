ElSimulator*
========

##Stack

* [Go](https://golang.org/)

**Feature Reader** Emulate server API (or other) with files.<br/>
**Feature Proxy** Helper to create file to read with proxy on the real API (or other).

```
./elSimulator -help
Usage of ./elSimulator:
  -baseDirectory="elSimulatorCurrent": directory with file to read (elSimulatorCurrent to use directory elSimulator)
  -bindingAddress="localhost:4000": The binding address
  -parameterRegex=".*": Parameter regex
  -proxyAddress="http://localhost:4000/file": The proxy address
```

##Start server ./elSimulator

####Feature Reader

One URL is one file [folder]/[name file].[detect ext] in directory **elSimulatorCurrent**.
We can overwritte http code and http headers with file [folder]**info_**[name file]**.json**. 

####Optional file [folder]info_[name file].json
```json
{
   "HttpCode": 200,
   "UrlRedirection": "",
   "Header": {
      "Header One": "Value One",
	  "Header Two": "Value Two"
   }
}
```

|URL|file|file info|http code|
|----|----|----|----|
|http://localhost:4000/file/test|current base/file/test/GET/withoutParameter.xml|current base/file/test/GET/info_withoutParameter.json|500|
|http://localhost:4000/file/test?param=value|current base/file/test/GET/param_value.json|current base/file/test/GET/info_param_value.json|203|


###Feature Proxy (proxyAddress=http://www.google.fr)

|URL|called url|file|
|----|----|----|
|http://localhost:4000/proxy/?q=test|http://www.google.fr/?q=test|current base/proxy/GET/q_test|

**Inspired by work at M6 and BT*
