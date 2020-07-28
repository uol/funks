package funks_test

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/uol/funks"
	gthttp "github.com/uol/gotest/http"
	utils "github.com/uol/gotest/utils"
)

/**
* The util/collections library tests.
* @author rnojiri
**/

// ConfigDuration - an example to parse
type ConfigDuration struct {
	SomeDuration funks.Duration `json:"someDuration"`
}

// TestSyncMapSize - tests the function
func TestSyncMapSize(t *testing.T) {

	m := sync.Map{}

	assert.Equal(t, 0, funks.GetSyncMapSize(&m), "expected 0")

	m.Store("a", 1)

	assert.Equal(t, 1, funks.GetSyncMapSize(&m), "expected 1")

	m.Store("b", 10)

	assert.Equal(t, 2, funks.GetSyncMapSize(&m), "expected 2")

	for i := 0; i < 100; i++ {
		m.Store(strconv.Itoa(i), i)
	}

	assert.Equal(t, 102, funks.GetSyncMapSize(&m), "expected 102")

	m.Delete("b")

	assert.Equal(t, 101, funks.GetSyncMapSize(&m), "expected 101")
}

// TestTOMLDurationParse - tests the toml duration parse for configurations
func TestTOMLDurationParse(t *testing.T) {

	strDuration := fmt.Sprintf("%ds", rand.Int63n(59))
	strConf := fmt.Sprintf("SomeDuration = \"%s\"", strDuration)

	c := &ConfigDuration{}

	_, err := toml.Decode(strConf, c)
	if !assert.NoError(t, err, "unexpected error parsing toml string") {
		return
	}

	assert.Equal(t, strDuration, c.SomeDuration.String())
}

// TestJSONDurationUnmarshal - tests the json unmarshal parse for configurations
func TestJSONDurationUnmarshal(t *testing.T) {

	strDuration := fmt.Sprintf("%ds", rand.Int63n(59))
	strConf := fmt.Sprintf(`{ "someDuration": "%s" }`, strDuration)

	c := &ConfigDuration{}
	err := json.Unmarshal([]byte(strConf), c)
	if !assert.NoError(t, err, "unexpected error parsing json string") {
		return
	}

	assert.Equal(t, strDuration, c.SomeDuration.String())
}

// TestJSONDurationMarshal - tests the json duration marshal for configurations
func TestJSONDurationMarshal(t *testing.T) {

	seconds := rand.Int63n(59)

	c := &ConfigDuration{
		SomeDuration: funks.Duration{Duration: time.Duration(seconds) * time.Second},
	}

	result, err := json.Marshal(c)
	if !assert.NoError(t, err, "unexpected error parsing json string") {
		return
	}

	assert.Equal(t, fmt.Sprintf(`{"someDuration":"%ds"}`, seconds), (string)(result))
}

// TestNewDurationString - tests creating a new instance using a string
func TestNewDurationString(t *testing.T) {

	_, err := funks.NewStringDuration("x")
	if !assert.Error(t, err, "unexpected parsing error") {
		return
	}

	seconds := rand.Int63n(59)
	strDuration := fmt.Sprintf("%ds", seconds)

	d, err := funks.NewStringDuration(strDuration)
	if !assert.NoError(t, err, "unexpected error parsing json string") {
		return
	}

	assert.Equal(t, utils.MustParseDuration(strDuration), d.Duration, "must be equal")
}

// TestNewDuration - tests creating a new instance using a time.Duration
func TestNewDuration(t *testing.T) {

	randomDuration := time.Duration(utils.RandomInt(1, 60)) * time.Second
	d := funks.NewDuration(randomDuration)

	assert.Equal(t, randomDuration, d.Duration, "must be equal")
}

// TestForceNewStringDurationPanic - tests creating a new instance using the forced method
func TestForceNewStringDurationPanic(t *testing.T) {

	defer func() {
		r := recover()
		if !assert.NotNil(t, r, "not nil error expected") {
			return
		}

		assert.Error(t, r.(error), "expected an error recovery")
	}()

	funks.ForceNewStringDuration("x")

	assert.FailNow(t, "the code must not reach this point")
}

// TestForceNewStringDuration - tests creating a new forced instance using a string
func TestForceNewStringDuration(t *testing.T) {

	seconds := rand.Int63n(59)
	strDuration := fmt.Sprintf("%ds", seconds)

	d := funks.ForceNewStringDuration(strDuration)

	assert.Equal(t, utils.MustParseDuration(strDuration), d.Duration, "must be equal")
}

// TestHTTPClientConfiguration - tests a http client
func TestHTTPClientConfiguration(t *testing.T) {

	randomWait := time.Duration(utils.RandomInt(1, 3)) * time.Second

	serverConf := &gthttp.Configuration{
		Host:        "localhost",
		Port:        18080,
		ChannelSize: 5,
		Responses: map[string][]gthttp.ResponseData{
			"normal": {
				{
					RequestData: gthttp.RequestData{
						Method: "GET",
						URI:    "/test",
					},
					Status: http.StatusOK,
					Wait:   randomWait,
				},
			},
		},
	}

	s := gthttp.NewServer(serverConf)
	defer s.Close()

	// generates a timeout
	client := funks.CreateHTTPClient(randomWait-time.Second, true)
	_, err := client.Get("http://localhost:18080/test")
	if !assert.Error(t, err, "expected an error") {
		return
	}

	// a normal request
	client = funks.CreateHTTPClient(randomWait+time.Second, true)
	resp, err := client.Get("http://localhost:18080/test")
	if !assert.NoError(t, err, "expected no error now") {
		return
	}

	if !assert.Equal(t, http.StatusOK, resp.StatusCode, "expected a 200 code") {
		return
	}

	// a limited request server
	client = funks.CreateHTTPClientAdv(randomWait+time.Second, true, 1)
	startTime := time.Now()
	for i := 0; i < 3; i++ {
		resp, err = client.Get("http://localhost:18080/test")
		if !assert.NoError(t, err, "expected no error now") {
			return
		}
	}
	elapsedTime := time.Since(startTime)

	assert.Equal(t, randomWait.Seconds()*3, math.Floor(elapsedTime.Seconds()), "expected only one request per %s", randomWait)
}
