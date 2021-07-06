# auto-test
Auto testing microservice for Authena project

## Description
This service checks attempts to pass testing tasks

### Input
* Test points — array of test points, point includes variants and correct variant numbers
* Grade scale — relationship between lower percentage and grade

### Output
* Percent — final test passing percentage
* Grade — final grade depending on the test passing percentage

### Microservice's adapter
This service uses gRPC as adapter