# Auction House

## This code provides a solution for a Programming Test


[![Code quality status](https://github.com/senseyman/auction-house-algorithm/actions/workflows/test-on-push.yml/badge.svg)](https://github.com/senseyman/auction-house-algorithm/actions/workflows/test-on-push.yml)


## Requirements
You can read the task requirement [here](./requirements.md)

To run this code, you need to install Golang and make.
* How to install [Golang](https://go.dev/doc/install)
* How to install [Make](https://www.gnu.org/software/make/)

## Before running the app
Before running the app, download the dependencies
```shell
make dep
```

You can also check if app tests passes
```shell
make test 
```

Also you can check code by go standards using the following command
```shell
make vet lint
```

## Run the app
Before starting the app, please prepare your test input file in text format.
Example of the file **input.txt**
```text
10|1|SELL|phone|10.00|20
12|8|BID|phone|7.50
13|5|BID|phone|12.50
15|8|SELL|laptop|250.00|20
16
17|8|BID|phone|20.00
18|1|BID|laptop|150.00
19|3|BID|laptop|200.00
20
21|3|BID|laptop|300.00
```

To run the app, execute the following command in the project root directory
```shell
go run main.go --path=input.txt
```
where *--path* takes the path to the input file with test execute instructions.

To get more information about input parameters, please run
```shell
go run main.go --help
```

The output of the app will be represented in stdout like the next example:
```text
20|phone|8|SOLD|12.50|3|20.00|7.50
20|laptop||UNSOLD|0.00|2|200.00|150.00
```