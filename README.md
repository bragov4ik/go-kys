[![Actions Status](https://github.com/bragov4ik/go-kys/workflows/build/badge.svg)](https://github.com/bragov4ik/go-kys/actions)
[![codecov](https://codecov.io/gh/bragov4ik/go-kys/branch/master/graph/badge.svg)](https://codecov.io/gh/bragov4ik/go-kys)

# Go knowledge yield summary

<img alt="asciicast" src="https://asciinema.org/a/gF52rkEfwDs1TNpK0d8PCNQqq.svg" width="600"/>

## Table of Content

- [Table of content](#table-of-content)
- [Glossary](#glossary)
- [Project Description](#project-description)
- [Importance](#importance)
- [Features](#features)
- [Installation](#installation)
- [Usage](#cli-usage)
- [How it works?](#how-it-works)
  - [Cyclomatic Complexity](#cyclomatic-complexity)
  - [Halstead Complexity](#halstead-complexity)
  - [Comments Complexity](#comments-complexity)
  - [Code Structure Complexity](#code-structure-complexity)
  - [Arithmetic Intricacy](#arithmetic-intricacy)
  - [Inline Data](#inline-data)
  - [Summing up](#summing-up)
- [Important Note](#important-note)
- [RUP filled template](#rup-filled-template)
- [Authors](#authors)
## Glossary
Weighted Micro Function Points ([WMFP](https://en.wikipedia.org/wiki/Weighted_Micro_Function_Points)) - a modern software sizing algorithm

## Project description
Evaluation of WMFP metrics in human-minutes for given Golang code. Implemented in Go.

Project developed for the course Software Systems Analysis and Design (SSAD) at Innopolis University in F21 semester.

## Importance
Software metrics help to estimate the size, price, or time consumption of software. There are no known open source solutions for counting the WMFP metrics, especially for Go. Thus, having a free alternative to proprietary software can be in demand.

## Features
* Count WMFP metric for a given golang source file
* Ability to choose multiple files explicitly to count total metrics for them
* Possibility to pick a folder and have all matching files summarized in metrics score
* Recursive mode of scanning a folder

## Installation
```console
$ go install github.com/bragov4ik/go-kys/cmd/gokys@latest
```

## CLI Usage
Default path to config file is `config.xml`
```console
$ cd $HOME/go/bin
$ ./gokys -c <PATH_TO_CONFIG> a.go           # calculates for file
$ ./gokys -c <PATH_TO_CONFIG> .              # calculates for project
$ ./gokys -c <PATH_TO_CONFIG> a.go b.go c.go # calculates for multiple files
```
## How it works?
### Cyclomatic Complexity

The [Cyclomatic Complexity](https://en.wikipedia.org/wiki/Cyclomatic_complexity)
indicates the complexity of a program.

This program calculates the complexities of each function by counting independent paths. It starts with the initial value
of 1 and each time program encounters one of the `if, for, case, ||, &&` statements it increases the value by
corresponding to the statement's weight specified in the configuration file.

### Halstead Complexity
Calculation of Halstead Complexity can be found [here](https://en.wikipedia.org/wiki/Halstead_complexity_measures)

### Comments Complexity
Measures the amount of effort spent on writing program comments. It calculates the number of words written in comments
and multiply to the word's weight specified in the configuration file

### Code Structure Complexity
Measures the amount of effort spent on the program structure such as separating code into classes, functions, and
interfaces. It starts with the initial value of 0 and each time program encounters structure declaration or function
declaration or interface declaration it increases value by declaration's weight specified in the configuration file

### Arithmetic Intricacy
Measures the complexity of arithmetic calculations across the program. It starts with initial value of 0 and each time
program encounters one of the `+ - * / % += -= *= /= %= ++ --` operators it increases value by operator's weight
specified in configuration file

### Inline Data
Measures the amount of effort spent on embedding hard-coded data. It starts with the initial value of 0 and each time
program encounters basic literal or composite literal it increases value by literal's weight specified in the configuration
file

### Summing Up
A program sums up all the above metrics to calculate total effort in human-minutes.

## Contribution
To contribute to the project fork the repository and make a PR.

## Important Note
Since there is no comprehensive description of the work of WMFP in open sources, some metrics were not implemented and some
were implemented using our vision and understanding of these metrics.

## RUP filled template
[Link](https://docs.google.com/document/d/1su-LKhZ33DbZ898iwvInVrTbZTy12idO/edit?usp=sharing&ouid=106194539643127537689&rtpof=true&sd=true)

## Authors
* Kirill Ivanov
* Anatoliy Baskakov
* Iskander Bakhtiyarov
* Ivan Rybin
