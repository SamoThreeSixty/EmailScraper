package imap

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func FetchLastUnseen(c *client.Client, inbox string, limit int) ([]*imap.Message, error) {
	_, err := c.Select(inbox, false)
	if err != nil {
		return nil, err
	}

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}

	ids, err := c.Search(criteria)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, nil
	}

	if len(ids) > limit {
		ids = ids[len(ids)-limit:]
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	var result []*imap.Message
	for msg := range messages {
		result = append(result, msg)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return result, nil
}
