package kafkax

type KafkaConf struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Brokers  []string `json:"brokers"`
	GroupID  string   `json:"group_id"`
	Topic    string   `json:"topic"`
}
