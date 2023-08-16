## go-trace-plugins
项目说明：go-trace-plugins是一个基于go语言的trace插件集合。目前支持的插件有：gorm、go-redis、http、trace客户端注册。

## 包介绍
| 包名        | 介绍                                                       |
|-----------|----------------------------------------------------------|
| gormotel  | 该插件主要功能是在程序执行sql语句时生成trace span，并将sql信息保存到span内供链路追踪时查看。 |
| redisotel | 针对redis操作命令生成span信息。                                     |
| httpotel  | 针对http请求生成span信息。                                        |
| trace     | trace客户端注册包，用于注册trace客户端。目前支持jaeger、zipkin、otlp-http。    |

具体使用方法请参考包内的test文件。