package pubsub

var (
	// TopicRecordingBegin is Pub/Sub topic name that publishes when recording is began
	TopicRecordingBegin = "nejireco/recording/begin"

	// TopicRecordingEnd is Pub/Sub topic name that publishes when recording is ended
	TopicRecordingEnd = "nejireco/recording/end"
)

var topics = []string{
	TopicRecordingBegin,
	TopicRecordingEnd,
}
