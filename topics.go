package pubsub

var (
	// TopicRecordingBegin is Nejireco Pub/Sub topic id that publishes when recording is began.
	TopicRecordingBegin = "nejireco/recording/begin"

	// TopicRecordingEnd is Nejireco Pub/Sub topic id that publishes when recording is ended.
	TopicRecordingEnd = "nejireco/recording/end"
)

var topics = []string{
	TopicRecordingBegin,
	TopicRecordingEnd,
}
