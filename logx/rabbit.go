package logx

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

// LogMessage 定义日志消息结构
type LogMessage struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// RabbitMQWriter 实现 zapcore.WriteSyncer 接口
type RabbitMQWriter struct {
	channel *amqp.Channel
	queue   string // 队列名称
}

// NewRabbitMQWriter 创建一个新的 RabbitMQWriter
func NewRabbitMQWriter(amqpURL, queueName string) (*RabbitMQWriter, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// 声明队列
	_, err = ch.QueueDeclare(
		queueName, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 独占
		false,     // 不等待
		nil,       // 额外参数
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	return &RabbitMQWriter{
		channel: ch,
		queue:   queueName,
	}, nil
}

// Write 将日志写入 RabbitMQ
func (w *RabbitMQWriter) Write(p []byte) (n int, err error) {
	// 解析 JSON 日志
	var logMsg LogMessage
	err = json.Unmarshal(p, &logMsg)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal log message: %v", err)
	}

	// 重新编码为 JSON 发送到 RabbitMQ
	body, err := json.Marshal(logMsg)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal log message: %v", err)
	}

	err = w.channel.Publish(
		"",      // 默认交换机
		w.queue, // 路由键为队列名称
		false,   // 不强制
		false,   // 不延迟
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to publish message: %v", err)
	}

	return len(p), nil
}

// Sync 实现 WriteSyncer 接口（RabbitMQ 不需要同步）
func (w *RabbitMQWriter) Sync() error {
	return nil
}

// Close 关闭 RabbitMQ 连接
func (w *RabbitMQWriter) Close() error {
	return w.channel.Close()
}
