package moody

var (
	// TopicRecordingBegin is Nejireco Pub/Sub topic id that publishes when recording is began.
	TopicRecordingBegin = "nejireco/recording/begin"

	// TopicRecordingEnd is Nejireco Pub/Sub topic id that publishes when recording is ended.
	TopicRecordingEnd = "nejireco/recording/end"
)

// Topics are topics of Nejireco.
var Topics = []string{
	TopicRecordingBegin,
	TopicRecordingEnd,
}
