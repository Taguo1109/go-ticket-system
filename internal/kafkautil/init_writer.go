package kafkautil

/**
 * @File: init_writer.go.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/18 下午9:29
 * @Software: GoLand
 * @Version:  1.0
 */
import (
	"context"
	"encoding/json"
	"errors"
	"github.com/taguo1109/go-ticket-system/internal/model"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

var KafkaWriter *kafka.Writer

func InitWriter() {

	// 自動建立 topic（只會執行一次）
	EnsureTopic("ticket-booking", "localhost:9092", 1)

	KafkaWriter = &kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"),
		Topic:        "ticket-booking",
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 1 * time.Millisecond, // 降低等待時間
		RequiredAcks: kafka.RequireOne,     // 要求至少有一個 broker 確認送達
		Async:        false,                // 確保發送成功回應才繼續
	}
	log.Println("Kafka Writer 初始化完成")
}

func SendTicketRequest(req model.TicketRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Key:   []byte(req.UserID),
		Value: data,
	}
	return KafkaWriter.WriteMessages(context.Background(), msg)
}

func EnsureTopic(topic string, brokerAddr string, partitions int) {
	// 第一次連線到 Kafka broker（取得 controller）
	conn, err := kafka.Dial("tcp", brokerAddr)
	if err != nil {
		log.Fatalf("Kafka 連線失敗: %v", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Fatalf("取得 controller 錯誤: %v", err)
	}

	// 再連一次到 controller broker（才有建立 topic 權限）
	conn, err = kafka.Dial("tcp", controller.Host+":"+strconv.Itoa(controller.Port))
	if err != nil {
		log.Fatalf("連接 controller 失敗: %v", err)
	}
	defer conn.Close()

	// 嘗試建立 topic
	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: 1,
	})

	if err != nil {
		log.Printf("建立 topic [%s] 可能已存在: %v", topic, err)
	} else {
		log.Printf("Kafka topic [%s] 建立完成", topic)
	}
}

func BootstrapKafka(reader *kafka.Reader) error {

	log.Println("🚀 Bootstrapping Kafka Consumer...")

	start := time.Now()

	for {
		if time.Since(start) > 5*time.Second {
			return errors.New("❌ Kafka bootstrap timeout")
		}

		// 用短 timeout 試著拉一次資料
		subCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		_, err := reader.ReadMessage(subCtx)
		if err == nil {
			log.Println("✅ Kafka Ready (got message!)")
			return nil
		}
		if strings.Contains(err.Error(), "Request Timed Out") ||
			strings.Contains(err.Error(), "context deadline exceeded") {
			log.Println("⏳ 等待 Kafka Ready...")
			time.Sleep(200 * time.Millisecond)
			continue
		}
		return err
	}
}
