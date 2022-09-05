package check

import (
	"io"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
** TO RUN THIS TEST SUITE CLI SHOULD BE BUILD TO BE UP-TO-DATE BINARY.
** USE make build-cli FOR IT
 */

//mock keys for testing purposes
var mockKeys = []string{
	"3db24f06-bd6f-4baf-a9fc-8174baff90d5",
	"ca5a5000-3737-431c-9c43-03e8244008cd",
}

var testDataPath = "test_data"

//restoreKeys takes up mock-keys, truncates file to 0 bytes and writes on each line $key + \r\n, as API WriteKey does
func restoreKeys() {

	f, err := os.OpenFile("test_data/keys.txt", os.O_TRUNC|os.O_RDWR, 0755)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	err = f.Truncate(0)
	if err != nil {
		panic(err)
	}

	for _, key := range mockKeys {
		towr := key + "\r\n"
		_, err := f.Write([]byte(towr))
		if err != nil {
			panic(err)
		}
	}

}

//clearAllKeys would trunc file to 0 bytes resulting in 0 keys present
func clearAllKeys() {

	f, err := os.OpenFile("test_data/keys.txt", os.O_TRUNC|os.O_RDWR, 0755)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	err = f.Truncate(0)
	if err != nil {
		panic(err)
	}
}

func TestGetFirstKeyOK(t *testing.T) {
	defer restoreKeys()

	checkService := NewCheckService(testDataPath)

	//Existing key is first key that is stored in test_data/keys.txt already for testing purposes
	existingKey := "3db24f06-bd6f-4baf-a9fc-8174baff90d5"

	key, err := checkService.GetFirstKey()
	require.NoError(t, err)
	require.Equal(t, existingKey, key)
	t.Logf("obtained key - %s\n", key)
}
func TestGetFirstKeyNoApiKeysLeft(t *testing.T) {
	defer restoreKeys()

	//Firstly, clear all keys
	clearAllKeys()

	checkService := NewCheckService(testDataPath)
	key, err := checkService.GetFirstKey()
	require.Error(t, err)
	require.Equal(t, "", key)
	require.Equal(t, NoApiKeysLeft.Error(), err.Error())
}

func TestRestoreKeyOK(t *testing.T) {
	/*
		This test-case is a multistep:
		 - Call checkService.RestoreKey(), expect no errors (keys are present after each test-run, see restoreKeys impl..)
		 - Call checkService.GetFirstKey(), should return the second key in mockKeys array
	*/

	defer restoreKeys()

	checkService := NewCheckService(testDataPath)

	err := checkService.RestoreKey()
	require.NoError(t, err)
	require.Nil(t, err)

	key, err := checkService.GetFirstKey()
	require.NoError(t, err)
	require.Nil(t, err)
	require.Equal(t, mockKeys[1], key)
}
func TestRestoreKeyNoApiKeysLeft(t *testing.T) {

	defer restoreKeys()

	clearAllKeys()

	checkService := NewCheckService(testDataPath)

	err := checkService.RestoreKey()
	require.Error(t, err)
	require.Equal(t, NoApiKeysLeft.Error(), err.Error())
	require.NotNil(t, err)

}

func TestCopyOK(t *testing.T) {

	w := httptest.NewRecorder()

	checkService := NewCheckService(testDataPath)

	err := checkService.Copy(w)
	require.NoError(t, err)

	respBytes, err := io.ReadAll(w.Body)
	require.NoError(t, err)

	clenh := w.Header().Get("Content-Length")

	clenAsInt, err := strconv.ParseInt(clenh, 10, 64)
	require.NoError(t, err)

	require.Equal(t, int(clenAsInt), len(respBytes))
	require.Equal(t, 0, int(clenAsInt)-len(respBytes))

	require.NotEqual(t, "", string(respBytes))
}
