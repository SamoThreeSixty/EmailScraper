package format

import "github.com/emersion/go-imap"

func FormatEmailAddressList(list []*imap.Address) string {
	var result string
	for i, addr := range list {
		if i > 0 {
			result += ", "
		}
		result += addr.MailboxName + "@" + addr.HostName
	}
	return result
}
