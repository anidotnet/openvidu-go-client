package openvidu

type OpenViduRole int

const (
	// Can subscribe to published Streams of other users
	SUBSCRIBER OpenViduRole = iota

	// SUBSCRIBER permissions + can publish their own Streams (call
	// Session.publish())
	PUBLISHER

	// SUBSCRIBER + PUBLISHER permissions + can force the unpublishing or
	// disconnection over a third-party Stream or Connection (call
	// Session.forceUnpublish() and
	// Session.forceDisconnect())
	MODERATOR
)
