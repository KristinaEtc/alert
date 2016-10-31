package alert

import (
	"github.com/mqu/go-notify"
)

//PopUpNotify shows notification
func PopUpNotify(id string) {
	log.Debug("Notify")
	notify.Init("TIME-OUT")
	hello := notify.NotificationNew("TimeOut", "alert", id)
	hello.Show()
}
