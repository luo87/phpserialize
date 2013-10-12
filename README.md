# phpserialize

serialize/unserialize data with golang.

base on [unserialize.js](https://github.com/kvz/phpjs/blob/master/functions/var/unserialize.js) and [go-simplejson](https://github.com/bitly/go-simplejson).

## Install

    $ go get github.com/luobailiang/phpserialize

## Use


	import "github.com/luobailiang/phpserialize"
    data, err := phpserialize.Unserialize(`a:2:{i:232;s:17:"2323 中文 dfsdf";s:3:"sdf";a:3:{i:0;i:2;i:1;i:3;i:2;i:4;}}`)
    data.Get("232").String()
