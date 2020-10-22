json_filter支持用SQL的方式查询json格式的数据。虽然jq也可以筛选数据，但语法实在太难记，每次学完过一段时间就忘，而SQL则是天天在用的，想忘都忘不了，于是就有了这个项目。

比如文件test.log的内容如下:

```
{"level":"info","msg":"open file error","ts":1602259198}
{"level":"error","msg":"open file error","ts":1602259199}
{"level":"info","msg":"this is a string","ts":1602259200}
{"level":"error","msg":"close file error","ts":1602259201}
{"level":"info","msg":"foobar","ts":1602259202}
{"level":"info","msg":"foobar","ts":1602259203,"data":{"id":12345,"title":"hello","created_at":1602259100}}
```

当要查询level为error并且ts大于1602259199的数据时，可以这样:

```bash
cat test.log | json_filter -q "select * from t where level='error' and ts>1602259199"
```

结果为:

```
{"level":"error","msg":"close file error","ts":1602259201}
```

注意表名固定是t，不能是其他的名字。

如果只需要部分字段，也可以像写SQL那样取出指定字段:

``` bash
cat test.log | json_filter -q "select msg,data.title from t where data is not null"
```

结果为:

```
{"data.title":"hello","msg":"foobar"}
```

如果需要获取所有的key，还可以用一个特殊的字段名`[keys]`来获取:

```bash
cat test.log | json_filter -q "select [keys] from t where level='info'"
```

结果为:

```
{"[keys]":"level,msg,ts"}
{"[keys]":"level,msg,ts"}
{"[keys]":"level,msg,ts"}
{"[keys]":"data,level,msg,ts"}
```

目前支持的SQL关键字及运算符如下：

`(`、`)`、`+`、`-`、`*`、`/`、`%`、`=`、`>`、`<`、`>=`、`<=`、`<>`、`and`、`or`、`is null`、`is not null`、`like `、`not like`、`in`、`not in`
