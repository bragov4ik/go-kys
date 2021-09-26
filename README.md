# Go knowledge yield summary

## Project description
Project developed for the course Software Systems Analysis and Design (SSAD) at IU in F21 semester.

Evaluation of WMFP metrics for Golang code. It can be used for approximation of software size and measuring similar software development's the cost.

## Features
* Count WMFP metric for a given golang source file
* Ability to choose multiple files explicitly to count total metrics for them
* Possibility to pick a folder and have all matching files summarized in metrics score
* Recursive mode of scanning a folder
* Program can be installed using package managers

## CLI Usage
```console
$ ./gokys -f a.go # calculates for file
$ ./gokys -d . # calculates for project
$ ./gokys -f "a.go b.go c.go" # calculates for multiple files
```

## RUP filled template
[Link](https://docs.google.com/document/d/1su-LKhZ33DbZ898iwvInVrTbZTy12idO/edit?usp=sharing&ouid=106194539643127537689&rtpof=true&sd=true)

## Authors
* Kirill Ivanov
* Anatoliy Baskakov
* Iskander Bakhtiyarov
* Ivan Rybin
