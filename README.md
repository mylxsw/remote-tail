# RemoteTail

RemoteTail是一款支持同步显示多台远程服务器的日志文件内容更新的工具，使用它可以让你同时监控多台服务器中某个（某些）日志文件的变更，将多台服务器的`tail -f xxx.log`命令的输出合并展示。

![logo](https://oayrssjpa.qnssl.com/remote-tail.jpg)

## 使用场景

假设公司有两台web服务器A和B，由于初期没有专业运维进行配置集中式的日志服务系统，两台服务器上分别部署了两套相同的代码提供web服务，使用nginx作为负载均衡，请求根据设定的策略转发的这两台web服务器上。

AB两台服务器中的项目均将日志写到文件系统的`/home/data/logs/laravel.log`文件。这种情况下如果我们需要查看web日志是否正常，一般情况下就需要分别登陆两台服务器，然后分别执行`tail -f /home/data/logs/laravel.log`查看日志文件的最新内容，这在排查问题的时候是非常不方便的。
