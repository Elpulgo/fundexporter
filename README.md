# Fundexporter
Tool to print and export fund data from Avanza as csv-files.

## Download
+ Download and extract gzipped file from '/dist' folder for your operating system of choice.

## Configuration
+ Place a config.json file in the same folder where you choose to use your program (included in the gzipped folder).
+ Fundids can be found in the url on Avanza when you check out a fund.
+ Filter can currently be used with 3 different options
    + industry -> example: Ny teknik
    + region -> example: Sverige
    + index -> can only be used with the value 'Index', which indicates if index funds show be filtered or not.
+ Use as many filters as you want.
+ Results are sorted in descending order on 3M %.
```
{
    "fundIds": [
        11111,
        22222
    ],
    "filter": [
        "industry:aaa",
        "industry:bbb",
        "region:cccc",
        "region:ddd",
        "index:Index"
    ]
}
```

## How to use
```
$ ./fundexporter
$ fundexporter.exe

Example:

$ ./fundexporter report current
```
## Arguments
+ `No arguments passed`
    + If no argument is passed, top 10 funds of each filter will print.
+ `current`
    + All funds will be printed which are configured in the fundIds configuration.
+ `report`
    + A report, in csv format will be created, based on top 10 funds of each filter.
    + Will be saved in 'report.csv'
+ `report current`
    + A report, in csv format will be created, based on current funds.
    + Will be saved in 'report-current.csv'

## Csv formatting
- For current report the following csv-format is intended:

### Each row
_Empty for date (manual input)_ | _Bransch_ (manual input) | Id | Name | _Weight (manual input)_ | _Keep (manual input)_ | _Current value (manual input)_ | _Development currency(manual input)_ | _Development % (manual input)_ | NAV | OMXSPI | M1 % | M3 % | M6 % | M12 % | MA30 | MA50 | MA200 | RSI(14) 

### Below funds
_Empty for date (manual input)_ | UNRate description | Value

