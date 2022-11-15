package gologger

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	buffer := bytes.Buffer{}
	t.Run("WithField option", func(t *testing.T) {
		logger, err := NewLogger(WithOutput(&buffer), WithField("id", 0))
		assert := assert.New(t)
		assert.Nil(err)
		assert.NotNil(logger)
		logger.Info("test")
		fmt.Println("buffer: ", buffer.String())
		assert.True(strings.Contains(buffer.String(), "id=0"))
	})
	t.Run("WithOutPut option", func(t *testing.T) {
		logger, err := NewLogger(WithOutput(&buffer), WithField("id", 0))
		assert := assert.New(t)
		assert.Nil(err)
		assert.NotNil(logger)
		logger.Info("test")
		assert.True(strings.Contains(buffer.String(), "test"))
	})
	t.Run("Singleton logic", func(t *testing.T) {
		// buffer := bytes.Buffer{}
		assert := assert.New(t)
		wg := sync.WaitGroup{}
		wg.Add(2)
		for i := 1; i < 3; i++ {
			go func(i int) {
				// logger, err := NewLogger(WithOutput(&buffer))
				logger, err := NewLogger()
				assert.Nil(err)
				assert.NotNil(logger)
				// fmt.Printf("logger output: %v\n", logger.Logger.Out)
				// logger.WithField("id", i).Info("test")
				fmt.Println("inside goroutine: ", i)
				logger.Info("test")
				fmt.Println("buffer: ", buffer.String())
				wg.Done()
			}(i)
		}
		wg.Wait()
	})
}
