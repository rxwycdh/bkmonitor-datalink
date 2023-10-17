package notifier

import (
	"context"
	"crypto/sha512"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/internal/apm/pre_calculate/window"
	"github.com/xdg-go/scram"
	"time"
)

type KafkaConfig struct {
	KafkaTopic    string
	KafkaGroupId  string
	KafkaHost     string
	KafkaUsername string
	KafkaPassword string
}

func KafkaHost(h string) Option {
	return func(args *Options) {
		args.KafkaHost = h
	}
}

func KafkaUsername(u string) Option {
	return func(args *Options) {
		args.KafkaUsername = u
	}
}

func KafkaPassword(p string) Option {
	return func(args *Options) {
		args.KafkaPassword = p
	}
}

func KafkaTopic(t string) Option {
	return func(options *Options) {
		options.KafkaTopic = t
	}
}

func KafkaGroupId(g string) Option {
	return func(options *Options) {
		options.KafkaGroupId = g
	}
}

type kafkaNotifier struct {
	ctx context.Context

	config        KafkaConfig
	consumerGroup sarama.ConsumerGroup
	handler       consumeHandler
}

func (k *kafkaNotifier) Spans() <-chan []window.Span {
	return k.handler.spans
}

func (k *kafkaNotifier) Start() {
	logger.Infof("KafkaNotifier started. host: %s topic: %s groupId: %s", k.config.KafkaHost, k.config.KafkaTopic, k.config.KafkaGroupId)
loop:
	for {
		select {
		case <-k.ctx.Done():
			logger.Infof("ConsumerGroup stopped.")
			break loop
		default:
			if err := k.consumerGroup.Consume(k.ctx, []string{k.config.KafkaTopic}, k.handler); err != nil {
				logger.Errorf("ConsumerGroup fails to consume. error: %s", err)
				time.Sleep(1 * time.Second)
			}
		}
	}
}

type consumeHandler struct {
	spans chan []window.Span
}

func (c consumeHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c consumeHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c consumeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for {
		select {
		case msg := <-claim.Messages():
			session.MarkMessage(msg, "")
			var s window.OriginMessage
			if err := json.Unmarshal(msg.Value, &s); err != nil {
				logger.Warnf("an abnormal span format was received -> %s. error: %s", msg.Value, err)
				continue
			}
			if len(s.Items) != 0 {
				c.spans <- s.Items
			}
		case <-session.Context().Done():
			logger.Infof("kafka consume handler session done")
			return nil
		}
	}
}

func newKafkaNotifier(setters ...Option) Notifier {

	args := &Options{}

	for _, setter := range setters {
		setter(args)
	}
	kafkaConfig := args.KafkaConfig
	logger.Infof("listen %s topic as groupId: %s, establish a kafka[%s(%s:%s)] connection", kafkaConfig.KafkaTopic, kafkaConfig.KafkaGroupId, kafkaConfig.KafkaHost, kafkaConfig.KafkaUsername, kafkaConfig.KafkaPassword)

	authConfig := getConnectionSASLConfig(kafkaConfig.KafkaUsername, kafkaConfig.KafkaPassword)
	group, err := sarama.NewConsumerGroup([]string{kafkaConfig.KafkaHost}, kafkaConfig.KafkaGroupId, authConfig)
	if err != nil {
		logger.Errorf("Failed to create a consumer group, topic: %s may not be consumed correctly. error: %s", kafkaConfig.KafkaTopic, err)
	}
	return &kafkaNotifier{
		ctx:           args.ctx,
		config:        args.KafkaConfig,
		consumerGroup: group,
		handler:       consumeHandler{spans: make(chan []window.Span, args.chanBufferSize)},
	}

}

// getConnectionSASLConfig Establish a connection by SHA512
func getConnectionSASLConfig(username, password string) *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Version = sarama.V0_10_2_1

	if username != "" && password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &SCRAMSHA512Client{HashGeneratorFcn: SHA512} }
	}

	return config
}

var (
	SHA512 scram.HashGeneratorFcn = sha512.New
)

type SCRAMSHA512Client struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *SCRAMSHA512Client) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *SCRAMSHA512Client) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}
