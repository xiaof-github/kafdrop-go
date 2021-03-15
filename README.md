# kafdrop-go


## 1，关于kafdrop
做这个项目是为了学习golang 和 kafka。仿照kafdrop的部分功能实现了go的版本，实现功能：  
1.查看broker地址和controller。  
2.查看topic列表，每个topic分区个数，每个topic可消费消息个数。  
3.查看topic下前200个消息内容，消息offset，消息所在分区。  

![KafkaBroker](https://github.com/xiaof-github/kafdrop-go/blob/master/pic/70E4379A-A1C6-49cf-9353-CA87B0DB9748.png)
![KafkaTopic](https://github.com/xiaof-github/kafdrop-go/blob/master/pic/7B8BDB8F-7092-4ede-9AEC-A120F7C585C0.png)
![KafkaMessage](https://github.com/xiaof-github/kafdrop-go/blob/master/pic/82585CA5-71CF-4c3c-9E40-3FD8AC642043.png)


## 2，创建工程的过程
拷贝beego的样例代码，在此基础开发

## 3，页面的构建
http://v3.bootcss.com/
https://jquery.com/download/

## 4, 编译运行
go build -o kafdrop.exe main.go
在cmd窗口执行kafdrop.exe
访问http://127.0.0.1:8090

## 5, 在线demo
http://152.136.200.213:8090/
