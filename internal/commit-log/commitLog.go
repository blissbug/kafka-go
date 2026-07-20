package commit_log

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/yourusername/go-kafka/internal/segment"
)

type CommitLog struct {
	segments       []*segment.Segment
	activeSegment  *segment.Segment
	maxSegmentSize int64
}

func (c *CommitLog) Open() error {
	files, err := os.ReadDir("data")

	if err != nil {
		return err
	}

	n := 0
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".log" {
			files[n] = f
			n++
		}
	}

	files = files[:n]

	sort.Slice(files, func(i, j int) bool {
		firstBaseOffset, err := strconv.Atoi(files[i].Name()[:len(files[i].Name())-len(".log")])

		if err != nil {
			panic(err)
			return true
		}
		secondBaseOffset, err := strconv.Atoi(files[j].Name()[:len(files[j].Name())-len(".log")])

		if err != nil {
			panic(err)
			return true
		}

		return firstBaseOffset < secondBaseOffset
	})

	//sort by filename here
	for _, f := range files {
		fName := f.Name()

		if filepath.Ext(fName) != ".log" {
			continue
		}
		fileBaseOffset := fName[:len(fName)-len(".log")]
		fileBaseOffsetInt, err := strconv.ParseInt(fileBaseOffset, 10, 64)
		if err != nil {
			return err
		}
		s := segment.Segment{}
		err = s.Open(fileBaseOffsetInt, c.maxSegmentSize)
		if err != nil {
			return err
		}
		c.segments = append(c.segments, &s)
		c.activeSegment = &s
	}

	if c.segments == nil || len(c.segments) == 0 {
		//handle opening segments if present else create a new
		err := c.createNewSegment()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CommitLog) Append(data []byte) (pos int, err error) {
	if c.activeSegment.IsFull() {
		err = c.createNewSegment()
		if err != nil {
			return -1, err
		}
	}

	pos, err = c.activeSegment.Append(data)
	if err != nil {
		return -1, err
	}

	return pos, nil
}

func (c *CommitLog) Read(index int64) ([]byte, error) {
	for _, s := range c.segments {
		if s.BaseOffset <= index && index < s.NextOffset {
			data, err := s.Read(index)
			return data, err
		}
	}

	return nil, fmt.Errorf("index %d not found", index)
}

func (c *CommitLog) Close() error {
	for _, s := range c.segments {
		err := s.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CommitLog) createNewSegment() error {
	//create a new segment and assign as active segment
	newSegment := &segment.Segment{}
	var logicalOffset int64
	if c.activeSegment != nil {
		logicalOffset = c.activeSegment.NextOffset
	} else {
		logicalOffset = 0
	}
	err := newSegment.Open(logicalOffset, c.maxSegmentSize)
	if err != nil {
		return err
	}

	c.segments = append(c.segments, newSegment)
	c.activeSegment = newSegment
	return nil
}
