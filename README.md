# About ini [![Build Status](https://travis-ci.org/knq/ini.svg)](https://travis-ci.org/knq/ini) #

A simple [Go](http://www.golang.org/project/) package for manipulating 
[ini files](https://en.wikipedia.org/wiki/INI_file).

This package is mostly a simple wrapper around the [parser package](/parser)
also in this repository. The parser package was implemented by generating a 
[Pigeon](https://github.com/PuerkitoBio/pigeon/) parser from a
[PEG grammar](https://en.wikipedia.org/wiki/Parsing_expression_grammar).

## Why Another ini File Package? ##

In writing a semi-related package, I evaluated a number of existing ini
packages. The other packages did not have all features that were needed,
and did not work correctly in many cases. As such, it was necessary to write a
package that worked badly formatted ini files, and that provided a more robust
option for parsing.

## Installation ##

Install the package via the following:

    go get -u github.com/knq/ini

## Usage ##

The ini package can be used similarly to the following:

    package main

    import (
        "log"
        "github.com/knq/ini"
    )

    var (
        data = `
        firstkey = one

        [some section]
        key = blah ; comment

        [another section]
        key = blah`
    )

    func main() {
        f := ini.LoadString(data)
        s := f.Get("some section")
        log.Printf(">> %s\n", s.Get("key"))
        s.Set("key2", "another value")
        f.Write("out.ini")
    }