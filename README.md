# RemoteTail

RemoteTail是一款支持同步显示多台远程服务器的日志文件内容更新的工具，使用它可以让你同时监控多台服务器中某个（某些）日志文件的变更，将多台服务器的`tail -f xxx.log`命令的输出合并展示。相比于其他流行的日志手机工具，RemoteTail去掉了在远程服务器安装agent的必要，减小了与远程服务器的耦合，但需要注意的是，由于日志收集使用的是远程执行`tail`命令，因此如果进程退出重启后会出现日志重复或者丢失部分日志的风险。

RemoteTail只适应于简单的日志收集聚合，如果你不介意重启服务时日志丢失或者重复的问题，那么推荐你尝试一下。

![logo](https://ssl.aicode.cc/remote-tail.jpg?20161011)

## 使用场景

假设公司有两台web服务器A和B，由于初期没有专业运维进行配置集中式的日志服务系统，两台服务器上分别部署了两套相同的代码提供web服务，使用nginx作为负载均衡，请求根据设定的策略转发的这两台web服务器上。

AB两台服务器中的项目均将日志写到文件系统的`/home/data/logs/laravel.log`文件。这种情况下如果我们需要查看web日志是否正常，一般情况下就需要分别登陆两台服务器，然后分别执行`tail -f /home/data/logs/laravel.log`查看日志文件的最新内容，这在排查问题的时候是非常不方便的。RemoteTail就是为了解决这种问题的，开发人员可以使用它同步显示两台（多台）服务器的日志信息。

## 安装

在[release页面](https://github.com/mylxsw/remote-tail/releases)下载对应的`remote-tail-平台`可执行文件，将该文件加入到系统的`PATH`环境变量指定的目录中即可。

比如，Centos下可以放到`/usr/local/bin`目录。

    mv remote-tail-linux /usr/local/bin/remote-tail

## 使用方法

使用前需要宿主机建立与远程主机之间的[ssh公钥免密码登陆](https://aicode.cc/linux-mian-mi-ma-deng-lu.html)。

    remote-tail -hosts 'watcher@192.168.1.226,watcher@192.168.1.225' \
    -file '/usr/local/openresty/nginx/logs/access.log'

![demo](https://ssl.aicode.cc/remote-tail-demo.jpg?20161011)

> 如果服务器sshd监听的非默认端口22，可以使用`watcher@192.168.1.226:2222`这种方式指定其它端口。

### 简单的日志收集

日志聚合后作为单独文件存储，可以使用下面的方法

    nohup remote-tail -hosts 'watcher@192.168.1.226,watcher@192.168.1.225' -file '/usr/local/openresthy/nginx/logs/access.log' -slient=true > ./res.log &

> `-slient=true`参数用于指定RemoteTail不输出欢迎信息和控制台彩色字符，只输出纯净整洁的日志。

### 指定配置文件

通过使用`-conf`参数可以为命令指定读取的配置文件，配置文件为TOML格式，请参考`example.toml`文件。

配置文件`example.toml`：

    # 全局配置,所有的servers中tail_file配置的默认值
    tail_file="/data/logs/laravel.log"
    
    # tail 命令的选项，一般Linux服务器不需要设置此项，采用默认值即可
    # 如果是AIX等服务器，可能tail命令不支持下面这两个选项，可以修改该配置项为 "-f"
    #tail_flags="--retry --follow=name"

    # 服务器配置,可以配置多个
    # 如果不提供 password, 则默认使用系统配置的 ssh-agent 设置，
    # 你也可以通过指定 private_key_path 配置项来指定使用特定的私钥来登录 (private_key_path=/home/mylxsw/.ssh/id_rsa)
    # server_name, hostname, user 配置为必选,其它可选
    [servers]

    [servers.1]
    server_name="测试服务器1"
    hostname="test1.server.aicode.cc"
    user="root"
    tail_file="/var/log/messages"
    # 指定ssh端口，不指定的情况下使用默认值22
    port=2222

    [servers.2]
    server_name="测试服务器2"
    hostname="test2.server.aicode.cc"
    user="root"
    tail_file="/var/log/messages"
    tail_flags="-f"

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

你可以在github的issue中提出你的bug或者其它需求。

## Stargazers over time

[![Stargazers over time](https://starchart.cc/mylxsw/remote-tail.svg)](https://starchart.cc/mylxsw/remote-tail)
