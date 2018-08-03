package helpers

import ( 

  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/spf13/viper"
)

func TestLoadConfig (t *testing.T) {

	assert := assert.New(t)

	assert.Equal(LoadConfig(), nil)
}

