# RemoteTail

RemoteTail是一款支持同步显示多台远程服务器的日志文件内容更新的工具，使用它可以让你同时监控多台服务器中某个（某些）日志文件的变更，将多台服务器的`tail -f xxx.log`命令的输出合并展示。

![logo](https://oayrssjpa.qnssl.com/remote-tail.jpg)

## 使用场景

假设公司有两台web服务器A和B，由于初期没有专业运维进行配置集中式的日志服务系统，两台服务器上分别部署了两套相同的代码提供web服务，使用nginx作为负载均衡，请求根据设定的策略转发的这两台web服务器上。

AB两台服务器中的项目均将日志写到文件系统的`/home/data/logs/laravel.log`文件。这种情况下如果我们需要查看web日志是否正常，一般情况下就需要分别登陆两台服务器，然后分别执行`tail -f /home/data/logs/laravel.log`查看日志文件的最新内容，这在排查问题的时候是非常不方便的。RemoteTail就是为了解决这种问题的，开发人员可以使用它同步显示两台（多台）服务器的日志信息。

## 安装

下载项目`bin/`下对应的`remote-tail-平台`可执行文件，将该文件加入到系统的`PATH`环境变量指定的目录中即可。

比如，Centos下可以放到`/usr/local/bin`目录。

    mv remote-tail-linux /usr/local/bin/remote-tail

## 使用方法

使用前需要宿主机建立与远程主机之间的[ssh公钥免密码登陆](http://b.aicode.cc/linux/2015/04/27/Linux%E4%BD%BF%E7%94%A8SSH%E5%85%AC%E9%92%A5%E5%85%8D%E5%AF%86%E7%A0%81%E7%99%BB%E5%BD%95.html)。

    remote-tail -hosts 'watcher@192.168.1.226,watcher@192.168.1.225' \
    -file '/usr/local/openresty/nginx/logs/access.log'

![demo](https://oayrssjpa.qnssl.com/remote-tail-demo.jpg)

### 指定配置文件

通过使用`-conf`参数可以为命令指定读取的配置文件，配置文件为TOML格式，请参考`example.toml`文件。

配置文件`example.toml`：

    # 全局配置,所有的servers中tail_file配置的默认值
    tail_file="/data/logs/laravel.log"

    # 服务器配置,可以配置多个
    # 如果不提供password,则使用当前用户的ssh公钥,建议采用该方式,使用密码方式不安全
    # server_name, hostname, user 配置为必选,其它可选
    [servers]

    [servers.1]
    server_name="测试服务器1"
    hostname="test1.server.aicode.cc"
    user="root"
    tail_file="/var/log/messages"

    [servers.2]
    server_name="测试服务器2"
    hostname="test2.server.aicode.cc"
    user="root"
    tail_file="/var/log/messages"

    [servers.3]
    server_name="测试服务器3"
    hostname="test2.server.aicode.cc"
    user="demo"
    password="123456"

执行命令：

    remote-tail -conf=example.toml

## 如何贡献

欢迎贡献新的功能以及bug修复，**Fork**项目后修改代码，测试通过后提交**pull request**即可。

## 问题反馈

你可以在github的issue中提出你的bug或者其它需求，也可以通过以下方式直接联系我。

- 微博：[管宜尧](http://weibo.com/code404)
- 微信：mylxsw

![WEIXIN](https://oayrssjpa.qnssl.com/weixin.jpg)
