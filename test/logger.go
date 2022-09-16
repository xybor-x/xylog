package test

import "github.com/xybor-x/xylog"

// Event10Fields create an EventLogger with 10 fields like zap benchmark.
func Event10Fields(logger *xylog.Logger) *xylog.EventLogger {
	var elogger = logger.Event("event")
	elogger.Field("int", _tenInts[0])
	elogger.Field("ints", _tenInts)
	elogger.Field("string", _tenStrings[0])
	elogger.Field("strings", _tenStrings)
	elogger.Field("time", _tenTimes[0])
	elogger.Field("times", _tenTimes)
	elogger.Field("user1", _oneUser)
	elogger.Field("user2", _oneUser)
	elogger.Field("users", _tenUsers)
	elogger.Field("error", errExample)
	return elogger
}
