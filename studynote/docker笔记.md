docker镜像相关命令

docker images #查看本地镜像

docker search [image] #搜索镜像

hub.docker.com #docker官方镜像仓库
docker pull [image]:tag #拉取指定版本镜像

docker rmi [image] #删除镜像

docker容器相关命令

docker run [-i -t -d] -name=c1 [image]:tag /bin/bash#运行镜像 -i 一直运行 -t 进入交互模式 -name=c1 为容器命名 -d 后台运行

docker exec -it c1 /bin/bash #进入容器

exit #退出容器

docker ps [-a -q]#查看运行中的容器  -a 显示历史所有容器 -q 只显示容器ID
 
docker stop [container] #停止容器

docker start [container] #启动容器

docker rm [container] #删除容器
docker rm `docker ps -a -q` #删除所有容器
docker inspect [container] #查看容器详细信息

docker logs [container] #查看容器日志7

数据卷
概念:
数据卷是宿主机和容器之间共享数据的一种方式。

docker run -v /data(宿主机):/data(容器) [image]:tag #挂载数据卷

数据卷容器
docker run -v /data --name data-container [image]:tag #创建数据卷容器
















