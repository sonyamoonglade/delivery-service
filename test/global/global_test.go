package global_test

import (
	"testing"

	"github.com/sonyamoonglade/delivery-service/test/global"
	"github.com/stretchr/testify/require"
)

func TestInitLogger(t *testing.T) {

	logger := global.InitLogger()
	require.NotNil(t, logger)

}