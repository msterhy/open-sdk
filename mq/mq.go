package mq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// todo:给下载站上mq
// todo promethus来做监控
// RabbitMQ 配置
const (
	MQURL      = "amqp://root:123456@localhost:56720/"
	QUEUE_NAME = "user_updates"
)

// 初始化 RabbitMQ 连接
func setupRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(MQURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// 声明队列
	_, err = ch.QueueDeclare(
		QUEUE_NAME, // 队列名称
		true,       // durable
		false,      // auto-delete
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return conn, ch, nil
}

// 发送消息到队列
func sendMessage(ch *amqp.Channel, user User) error {
	body, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	err = ch.Publish(
		"",         // 默认交换机
		QUEUE_NAME, // 队列名称（Routing Key）
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Printf("Message sent: %s\n", string(body))
	return nil
}
