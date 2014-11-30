ElSimulator*
========

**Feature Reader** Emulate server API (or other) with files.<br/>
**Feature Proxy** Helper to create file to read with proxy on the real API (or other).

```
./elSimulator -help
Usage of ./elSimulator:
  -baseDirectory="elSimulatorCurrent": directory with file to read (elSimulatorCurrent to use directory elSimulator)
  -bindingAddress="localhost:4000": The binding address
  -parameterRegex=".*": Parameter regex
  -proxyAddress="http://localhost:4000/file": The binding address
```

##Start server ./elSimulator

####Feature Reader

One URL is one file in directory **elSimulatorCurrent**.
We can overwritte http code with file [folder]info_[name file].json.

####Optional file [folder]info_[name file].json
```json
{
	"HttpCode" : 203
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

*inspired by work at M6 and BT
