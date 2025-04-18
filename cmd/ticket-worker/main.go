package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/taguo1109/go-ticket-system/internal/kafkautil"
	"log"
	"time"
)

/**
 * @File: main.go.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/18 下午9:12
 * @Software: GoLand
 * @Version:  1.0
 */

func main() {
	log.Println("Ticket Worker 啟動中...")

	// 建立 Kafka reader（也就是 consumer）
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"}, // Kafka broker 位址（通常是 Kafka container 的 port）
		Topic:       "ticket-booking",           // 要訂閱的 Topic 名稱（由 API 發出訊息）
		GroupID:     "ticket-group",             // Consumer Group 名稱（同 group 只會收到一份訊息）
		MinBytes:    1,                          // 最小讀取訊息大小（10KB）
		MaxBytes:    10e6,                       // 最大讀取訊息大小（10MB）
		MaxWait:     10 * time.Millisecond,      // 降到 10ms 馬上吐資料
		StartOffset: kafka.FirstOffset,          // 強制從頭開始讀
	})
	// 2. 執行 BootstrapKafka → 等 Kafka 準備好
	if err := kafkautil.BootstrapKafka(reader); err != nil {
		log.Fatalf("Kafka 初始化失敗: %v", err)
	}
	// 持續監聽 Kafka 訊息
	for {
		// 從 Kafka 讀取一則訊息
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Kafka 讀取失敗: %v", err)
			continue
		}

		// 印出收到的訊息內容
		log.Printf("收到搶票請求: %s", string(m.Value))

		// ⚠️ 後續你可以在這裡加：
		// 1. 將 JSON 解析成結構（ticketRequest）
		// 2. Redis 查詢該區域是否還有票
		// 3. 若有票則扣減，否則回傳搶票失敗
		// 4. 將結果寫入 MySQL / 傳送到 order-create topic
	}
}
