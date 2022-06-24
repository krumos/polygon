package main

import (
    "io/ioutil"
 
    "gopkg.in/yaml.v2"
)

var Texts map[interface{}]string


func getText(fileName string) () {
    data, _ := ioutil.ReadFile(fileName)
    Texts = make(map[interface{}]string)
    yaml.Unmarshal([]byte(data), &Texts)
}
