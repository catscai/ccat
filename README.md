# ccat

基于Golang的TCP服务器框架

#### 使用例子

```go    
type EchoMessage struct {
    msg.BaseMessage
    name string
}

func (em *EchoMessage) Unpack(data []byte) error {
    em.name = string(data)
    return nil
}

func (em *EchoMessage) Pack() ([]byte, error) {
    return []byte(em.name), nil
}

// Deal 用户自己负责封包,回复
func Deal(ctx *iface.CatContext, request iface.IRequest, message imsg.IMessage) error {
    req := message.(*EchoMessage)
    // 创建回包
    echo := EchoMessage{
        name: "ccat",
    }
    echoData, _ := echo.Pack()
    header := msg.DefaultHeader{
        PackType:  TestPackType1RS,
        SessionID: request.GetHeaderPack().GetSessionID().(uint64),
        Data:      echoData,
    }
    // 发送
    ctx.SendMsg(&header)
    return nil
}

// DealSimpleMessage 不需要自己封包,解包和发送, 函数返回后自动发送rsp
func DealSimpleMessage(ctx *iface.CatContext, reqMsg, rspMsg imsg.IMessage) error {
    req := reqMsg.(*EchoMessage)
    rsp := rspMsg.(*EchoMessage)
    
    rsp.name = req.name
    return nil
}

// DealSimplePB pb类型消息, 不需要自己封包,解包和发送, 函数返回后自动发送rsp
func DealSimplePB(ctx *iface.CatContext, reqMsg, rspMsg proto.Message) error {
    req := reqMsg.(*test.TestRQ)
    rsp := rspMsg.(*test.TestRS)
    rsp.Reply = proto.String("a reply message")
    ctx.Info("DealSimplePB recv data", zap.Any("req", *req), zap.Any("rsp", *rsp))
    return nil
}

// 包类型
const (
	TestPackType1RQ uint32 = 1
	TestPackType1RS uint32 = 2

	TestPackType2RQ uint32 = 3
	TestPackType2RS uint32 = 4

	TestPackType3RQ uint32 = 5
	TestPackType4RS uint32 = 6
)

func main() {
    service := ccat.GetServer("tcp_test")
    // 注册消息
    service.GetDispatcher().RegisterHandler(TestPackType1RQ, &EchoMessage{}, Deal)
    // 注册消息,处理回调之后, 自动回包
    service.GetDispatcher().RegisterHandlerSimple(TestPackType2RQ, TestPackType2RS, &EchoMessage{}, &EchoMessage{}, DealSimpleMessage)
    // 注册消息pb类型,处理回调之后, 自动回包
    service.GetDispatcher().RegisterHandlerSimplePB(TestPackType3RQ, TestPackType3RS, &test.TestRQ{}, &test.TestRS{}, DealSimplePB)
    ccat.Run()
}
```

### Server

##### 简介

+ 消息包结构为:[packLen][header][data]
+ 收到消息后的处理是交给工作者池处理的。工作池的参数可在配置文件中设置
+ 日志使用zap日志库简单封装

##### 关键接口

+ SetDataPack 自定义处理TCP粘包
+ SetDispatcher 自定义消息分发器
+ SetHeaderOperator 自定义包头解析规则

+ GetWorkerGroup().SetShardWorkerFunc 自定义消息分片规则， 如：我们需要跟根据包中的UID属性将用户消息排队, 控制并发,那么就可以设置处理worker的选择规则,workerID = UID %
  Size

### Connection

+ SendMsg 发送消息包
+ SendData 发送raw数据
+ SetProperty 在连接上携带属性
+ GetProperty 获取属性
+ RemoveProperty 删除属性

#### 配置文件

使用yaml格式配置

```yaml    
log:
    app_name: test-svr
    log_dir: ./logs/  # 默认./logs/
    level: debug      # 日志等级fatal/panic/dPanic/error/warn/info/debug
    max_size: 128     # 单个日志文件最大M
    max_age: 7        # 日志文本保存时间
    #max_backups: 30   # 项目日志最大副本数
    #func: false      # 是否打印函数名,默认false
    console: true    # 控制台是否打印,默认false
    
tcp_service:
    tcp_test: # 服务名称 用于区分多个service
        ip: 127.0.0.1   
        port: 2233
        max_conn: 1024  # 最大连接数,默认0-不限制
        auto: true  # 自动创建tcp服务并监听
        # max_pack_len: 65535 # 最大包长度,默认为0-不限制
        worker_group: # 工作者组配置
            size: 10        # 工作者数量
            queue_length: 10  # 包队列长度
```



