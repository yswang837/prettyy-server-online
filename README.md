# prettyy的服务端

* 概述：用 golang + docker 的方式搭建**企业级 web 后端**系统。所采用的技术主要有：golang、gin、mysql、redis、kafka、consul、docker、prometheus、granfna、zap日志等技术，以达到构建后端web服务的目的。

## 1、环境说明

* 环境差不多就行，我开发这个项目就是inter的mac和m1的mac混用的，要求不严格

1、masOS Big Sur 11.4 Apple M1

2、go版本 go1.19.5

3、docker 桌面版 server和client都是20.10.8

## 2、搭建mysql，redis集群，并配置主从同步

* docker network create --subnet=172.18.0.0/16 prettyy_net // 其ip范围是172.18.0.0-172.18.255.255 (只需要执行一次)
* mysql搭建，主从同步(用户名和密码都是root),conf在容器的/etc/mysql/my.cnf，https://blog.csdn.net/agonie201218/article/details/121499881
* docker run --name=prettyy-net-mysql-master --net prettyy_net --ip 172.18.0.10 --env=MYSQL_ROOT_PASSWORD=root
  --env=TZ=Asia/Shanghai --env=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin --env=GOSU_VERSION=1.14
  --env=MYSQL_MAJOR=8.0 --env=MYSQL_VERSION=8.0.28-1debian10 --volume=/Users/yuanshun/test-docker/mysql-master:/root
  --privileged -p 4445:3306 --restart=always mariadb:latest
* docker run --name=prettyy-net-mysql-slave --net prettyy_net --ip 172.18.0.11 --env=MYSQL_ROOT_PASSWORD=root
  --env=TZ=Asia/Shanghai --env=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin --env=GOSU_VERSION=1.14
  --env=MYSQL_MAJOR=8.0 --env=MYSQL_VERSION=8.0.28-1debian10 --volume=/Users/yuanshun/test-docker/mysql-slave:/root
  --privileged -p 4446:3306 --restart=always mariadb:latest
* redis搭建，主从同步，https://cloud.tencent.com/developer/article/1343837
* docker run -it --name prettyy-net-redis-master --net prettyy_net --ip 172.18.0.20 -d -p 4455:6379 redis redis-server 
* docker run -it --name prettyy-net-redis-slave --net prettyy_net --ip 172.18.0.21 -d -p 4456:6379 redis redis-server
* 然后进入slave，执行 slaveof 172.18.0.20 6379搞定

## 2、以docker方式打包项目，并启动项目

1、启动docker的守护进程

2、在项目根目录prettyy/server开启module模块管理，并且用该命令打包镜像 sudo docker build -f docker/app-builder -t prettyy-image-v0.0.1 .   (另一种方式是：不用dockerfile，创个go的容器，在容器中装git，拉取代码打包也行，本项目采取第一种)

3、docker run -it --name=prettyy-builder --net prettyy_net -v /Users/yuanshun/prettyy-docker:/root  -d prettyy-image-v0.0.1_id，然后进入容器，cp app /root

4、docker run --name=prettyy-server --ip 172.18.100.100 --net prettyy_net --env=GIN_MODE=release --env=HTTP_SERVER_LISTEN_ADDR=:7979 --env=BLOG_METRICS_SERVER_LISTEN_ADDR=:8989 --env=SERVICE_NAME=blog-server --env=IDC=aliyun  --ip 172.18.100.100 -p 7979:7979 -p 8989:8989 -p 9999:8500 -v /Users/yuanshun/prettyy-docker:/root -d docker_consul_id
* 进入prettyy-server容器/root 然后./app，服务就起来了，该服务可以访问mysql、redis等集群，本地请求localhost:7979/xxx就可以访问到服务，

## 3、服务名、ip、端口映射

1. prettyy-server       172.18.100.100    7979:7979    内置的 metrics 8989:8989，内置的 consul 9999:8500
2. prettyy-mysql-master 172.18.0.10       4445:3306
3. prettyy-mysql-slave  172.18.0.11       4446:3306
4. prettyy-redis-master 172.18.0.20       4455:6379
5. prettyy-redis-slave  172.18.0.21       4456:6379

## 3、环境变量说明

* GIN_MODE 默认是debug，线上需使用release
* HTTP_SERVER_LISTEN_ADDR 默认是系统随机寻找一个可用的端口，也可指定固定的端口,如 `:8080`
* BLOG_METRICS_SERVER_LISTEN_ADDR 默认是系统随机寻找一个可用的端口，也可指定固定的端口,如 `:9999`

## 5、项目wiki

**done**

1、创建github仓库、跑出hello world

2、引入gin框架、基于dockerfile创建镜像、部署hello world代码，将hello world跑到浏览器中、实现容器化部署

3、为方便，将第三方库都写入custom-pkg

5、dockerfile支持consul和golang

**todo**

6、生成user_0,user_1表,inverted_index_0,inverted_index_1表

```sql
CREATE TABLE `user_0` (
  `uid` int(11) NOT NULL DEFAULT 0 COMMENT 'user id',
  `email` varchar(64) NOT NULL DEFAULT '' COMMENT 'email',
  `password` varchar(32) NOT NULL DEFAULT '' COMMENT 'password',
  `phone` varchar(11) NOT NULL DEFAULT '' COMMENT 'phone number',
  `nick_name` varchar(32) NOT NULL DEFAULT '' COMMENT 'nick name',
  `role` int(11) NOT NULL DEFAULT 1 COMMENT 'role number',
  `grade` int(11) NOT NULL DEFAULT 1 COMMENT 'grade number',
  `avatar` varchar(64) NOT NULL DEFAULT '' COMMENT 'avatar',
  `summary` varchar(200) NOT NULL DEFAULT '' COMMENT 'summary',
  `gender` varchar(8) NOT NULL DEFAULT '' COMMENT 'gender',
  `province_city` varchar(64) NOT NULL DEFAULT '' COMMENT 'province_city',
  `code_age` int(11) NOT NULL DEFAULT 0 COMMENT 'code_age',
  `birthday` varchar(16) NOT NULL DEFAULT '' COMMENT 'birthday',
  `address` varchar(256) NOT NULL DEFAULT '' COMMENT 'address',
  `is_certified` int(11) NOT NULL DEFAULT 2 COMMENT 'is_certified',
  `data_integrity` int(11) NOT NULL DEFAULT 0 COMMENT 'data_integrity',
  `create_time` datetime NOT NULL DEFAULT current_timestamp() COMMENT 'create time',
  `update_time` datetime NOT NULL DEFAULT current_timestamp() COMMENT 'update time',
  `login_time` datetime NOT NULL DEFAULT current_timestamp() COMMENT 'login time',
  PRIMARY KEY (`uid`),
  UNIQUE (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPRESSED;

CREATE TABLE `inverted_index_0` (
  `id` int(11) NOT NULL NOT NULL AUTO_INCREMENT COMMENT 'id',
  `typ` varchar(8) NOT NULL DEFAULT '' COMMENT 'typ',
  `attr_value` varchar(64) NOT NULL DEFAULT '' COMMENT 'attr_value',
  `idx` varchar(64) NOT NULL DEFAULT '' COMMENT 'idx',
  `create_time` datetime NOT NULL DEFAULT current_timestamp() COMMENT 'create time',
  `update_time` datetime NOT NULL DEFAULT current_timestamp() COMMENT 'update time',
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPRESSED;

CREATE TABLE `article_0` (
   `aid` varchar(32) NOT NULL DEFAULT '' COMMENT 'article id',
   `title` varchar(128) NOT NULL DEFAULT '' COMMENT 'title',
   `content` longtext NOT NULL DEFAULT '' COMMENT 'content',
   `cover_img` varchar(256) NOT NULL DEFAULT '' COMMENT 'cover img',
   `summary` varchar(1024) NOT NULL DEFAULT '' COMMENT 'summary',
   `tags` varchar(256) NOT NULL DEFAULT '' COMMENT 'tags',
   `visibility` varchar(8) NOT NULL DEFAULT '' COMMENT 'visibility',
   `typ` varchar(8) NOT NULL DEFAULT '' COMMENT 'typ',
   `share_num` int(11) NOT NULL DEFAULT 0 COMMENT 'share num',
   `comment_num` int(11) NOT NULL DEFAULT 0 COMMENT 'comment num',
   `like_num` int(11) NOT NULL DEFAULT 0 COMMENT 'like num',
   `read_num` int(11) NOT NULL DEFAULT 0 COMMENT 'read num',
   `collect_num` int(11) NOT NULL DEFAULT 0 COMMENT 'collect num',
   `status` varchar(8) NOT NULL DEFAULT '' COMMENT 'status',
   `uid` int(11) NOT NULL DEFAULT 0 COMMENT 'user id',
   `create_time` datetime NOT NULL DEFAULT current_timestamp() COMMENT 'create time',
   `update_time` datetime NOT NULL DEFAULT current_timestamp() COMMENT 'update time',
    PRIMARY KEY (`aid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPRESSED;
CREATE TABLE `column_0` (
    `cid` varchar(32) NOT NULL DEFAULT '' COMMENT 'column id',
    `title` varchar(64) NOT NULL DEFAULT '' COMMENT 'title',
    `cover_img` varchar(256) NOT NULL DEFAULT '' COMMENT 'cover img',
    `summary` varchar(256) NOT NULL DEFAULT '' COMMENT 'summary',
    `front_display` varchar(8) NOT NULL DEFAULT '1' COMMENT 'front display',
    `is_free_column` varchar(8) NOT NULL DEFAULT '1' COMMENT 'is free column',
    `subscribe_num` int(11) NOT NULL DEFAULT 0 COMMENT 'subscribe num',
    `uid` int(11) NOT NULL DEFAULT 0 COMMENT 'user id',
    `create_time` datetime NOT NULL DEFAULT current_timestamp() COMMENT 'create time',
    `update_time` datetime NOT NULL DEFAULT current_timestamp() COMMENT 'update time',
    PRIMARY KEY (`cid`),
    UNIQUE (`title`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPRESSED;

```





7、封装auth中间件、注册用md5(pin+email+client_ip)，别的所有接口均采用jwt鉴权

8、实现mysql邮件注册功能，jwt接口鉴权模块，丐版

9、搭建容器网络的redis资源，并完成主从同步[](https://)

10、用cobra和gorm的自动迁移，生成mysql user数据库表，user表包含了以下字段

uid：主键，通过雪花算法生成，GenID()加个type参数AA代表用户,AB代表学校等

邮箱：至少需通过验证码实名注册一个，可通过这两个字段登录，前期用邮箱验证码或密码登录，后期完善短信验证码登录，两个都绑定，获得认证成就

电话：至少需通过验证码实名注册一个，可通过这两个字段登录，前期用邮箱验证码或密码登录，后期完善短信验证码登录，两个都绑定，获得认证成就

是否认证：目前取决于邮箱和电话，后期可增加姓名，身份证、住址等

昵称：拒绝非法昵称

密码：前期可用简单校验，后期完善复杂校验密码

头像：鼓励用户上传头像，若未上传，随机选择一个

角色：1 普通用户、2 管理员、3 超级管理员、1默认注册的角色类型，23不提供注册接口，联系系统管理员注册

等级：类似于qq等级，显示距下次升级还有多少天，用一个类似对数函数的方法来升级，等级越高升级越困难

码龄：通过注册时间来提供

荣誉墙：类似于qq的荣誉墙，比如连续登录成就，活跃发言成就，资料完成成就

创建时间：计算码龄等

更新时间：用于控制敏感字段的更新频率，如密码、手机号、邮箱、性别、身份证

个人简介：默认为：这个人很神秘~

性别：男女不填，一旦填了就无法修改

资料完整度：

外键：

## api管理
