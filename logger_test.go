package gologger

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	buffer := bytes.Buffer{}
	assert := assert.New(t)
	t.Run("WithOutPut option", func(t *testing.T) {
		logger, err := NewLogger(WithOutput(&buffer), WithField("id", 0))
		assert.Nil(err)
		assert.NotNil(logger)
		buffer.Reset()
		logger.Info("testing WithOutPut option")
		assert.True(strings.Contains(buffer.String(), "testing WithOutPut option"))
	})
	t.Run("WithField option", func(t *testing.T) {
		logger.Info("testing WithField option")
		assert.True(strings.Contains(buffer.String(), "id"))
	})
	t.Run("Singleton logic", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(2)
		for i := 1; i < 3; i++ {
			go func(i int) {
				newBuffer := bytes.Buffer{}
				newLogger, err := NewLogger(WithOutput(&newBuffer), WithField("id", i))
				assert.Nil(err)
				assert.NotNil(newLogger)
				newLogger.Info("testing Singleton logic")
				assert.False(strings.Contains(newBuffer.String(), "testing WithOutPut option"))
				assert.False(strings.Contains(newBuffer.String(), "0"))
				assert.True(strings.Contains(buffer.String(), "testing WithOutPut option"))
				assert.True(strings.Contains(buffer.String(), "0"))
				assert.False(strings.Contains(buffer.String(), "1"))
				wg.Done()
			}(i)
		}
		wg.Wait()
	})
}
