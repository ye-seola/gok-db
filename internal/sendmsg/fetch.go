package sendmsg

import (
	"fmt"
	"os"
	"regexp"
)

const perfPath = "/data_mirror/data_ce/null/0/com.kakao.talk/shared_prefs/KakaoTalk.hw.perferences.xml"

func mustFetchNotificationReferer() string {
	data, err := os.ReadFile(perfPath)
	if err != nil {
		panic(fmt.Errorf("mustFetchNotificationReferer: %w", err))
	}

	re := regexp.MustCompile(`<string name="NotificationReferer">([^<]+)</string>`)
	match := re.FindSubmatch(data)

	if len(match) > 1 {
		return string(match[1])
	} else {
		panic(fmt.Errorf("mustFetchNotificationReferer: failed to fetch data"))
	}
}
