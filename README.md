# 阿里云按量使用

由于阿里云ecs包月比较贵，而且我们并不是每时每刻都在使用，因此按量进行计费，需要时开机，不需要时关机，是最理想的使用情况
我们将会实现，点击启动或者关闭远程阿里云ecs实例

## 申请阿里云ecs

访问 [阿里云ECS购买控制台](https://ecs-buy.aliyun.com/ecs)

### ecs

1. 选择配置

![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606000128947.png)

**选择 按量付费 -> 地域 -> 2核2G或者2核4G**

![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606000617525.png)

**选择Ubantu -> 发行版本 -> 40GB云盘**

![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606001254043.png)

**设置登录密码 -> 下单**

由于ecs控制台权限更高，所以即使忘记密码也能通过控制台登录，这个是方便ssh连接
等待创建完成，然后登陆[ecs控制台](https://ecs.console.aliyun.com/server)

### 弹性网卡

由于我们关闭远程服务器再重启后，公网ip是会变化的，因此我们需要申请**弹性网卡**（类似于一个不变的公网ip）
1. 绑定公网ip
![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606002617365.png)

![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606002710166.png)

![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606074208327.png)

来到[弹性公网ip控制台](https://vpcnext.console.aliyun.com/eip)
将刚创建好的弹性网卡绑定到我们之前申请的ecs实例上

![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606075042923.png)
注意，这边显示的 **ip **，即是后续需要使用的 **&{host}**

![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606075407490.png)

### 登陆ecs，创建普通用户并赋予管理员权限

回到[阿里云ECS购买控制台](https://ecs-buy.aliyun.com/ecs)
远程连接
![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606091651134.png)
![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606080336332.png)
![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606080442287.png)

1. 创建新用户，系统会提示设置密码并填写用户信息（信息可直接按回车跳过）
!!! **注意这边将会设置后续我们需要使用到的\${name}和\${password}**
```bash
sudo adduser ${name}
```
2. 赋予管理员权限，用户添加到 sudo 组

```bash
sudo usermod -aG sudo ${name}
```

3. 切换至新用户并测试 sudo：

```bash
su -${name}
sudo whoami
```
![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606081717163.png)

出现root则说明设置成功

## 获取配置

我们一共需要获取以下配置信息，其中\${name} 和\${password} 我们刚刚设置
\${host}在设置弹性网卡得到

```yaml
server:
  host: ${host}
  name: ${name}
  password: ${password}

ali-ecs:
  endpoint: ${endpoint} # 阿里云结点
  instanceId: ${instanceId} # ecs实例id
  
  accessKeyId: ${accessKeyId} 
  accessKeySecret: ${accessKeySecret}

```

### 获取\${endpoint}和\${instanceId} .

访问 [ecs产品](https://api.aliyun.com/product/Ecs)
找到对应的区域
![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606083340169.png)
比如我是华南1（深圳），则\${endpoint} 为 ecs.cn-shenzhen.aliyuncs.com

回到[ecs控制台](https://ecs.console.aliyun.com/server)
![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606083717479.png)
获取对应的\${instanceId}

### 获取\${accessKeyId}和${accessKeySecret}.

访问[创建RAM用户](https://ram.console.aliyun.com/users)
点击 创建用户
![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606084615721.png)

![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606084815147.png)
这两个为我们需要的 \${accessKeyId}和${accessKeySecret}

为RAM用户授权ECS，回到[创建RAM用户](https://ram.console.aliyun.com/users)
![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606085215637.png)
![](https://imgstroe-redblacktree.oss-accelerate.aliyuncs.com/img/20260606092115359.png)


## 运行代码

### 修改根目录下的config.yaml.example

将config.yaml.example 重命名为 **config.yaml**
并修改里面对应配置占位符

### 编译出二进制文件（非必须，可选）

确保安装完成 go

```bash
# 安装依赖包
go mod tidy

# 编译二进制文件并重命名为remote-ssh
go build -o remote-ssh.exe
```

本身根目录下也有我已经编译好的 remote-ssh.exe文件，如果没有配置go,可以跳过这一步

### 使用说明

remote-ssh.exe 输入s：start  c：close
或者根目录下脚本：
start.bat  开机
close.bat 关机

