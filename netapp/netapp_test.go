package netapp_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/pepabo/go-netapp/netapp"
)

func setup() (baseURL string, mux *http.ServeMux, teardownFn func()) {
	mux = http.NewServeMux()
	srv := httptest.NewServer(mux)
	return srv.URL, mux, srv.Close
}

func fixture(path string, t *testing.T) []byte {
	r, err := os.OpenFile("fixtures/"+path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func createTestClientWithFixtures(t *testing.T) (c *netapp.Client, teardownFn func()) {
	baseURL, mux, teardown := setup()

	requestFixture := bytes.TrimSpace(fixture(fmt.Sprintf("%s_%s", t.Name(), "request.xml"), t))
	responseFixture := bytes.TrimSpace(fixture(fmt.Sprintf("%s_%s", t.Name(), "response.xml"), t))

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		bs, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("Got error reading body %s", err)
		}

		if !bytes.Equal(bs, requestFixture) {
			t.Errorf("%s: result not as expected:\n%v", t.Name(), diff.LineDiff(string(requestFixture), string(bs)))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(responseFixture)
	})

	c = netapp.NewClient(baseURL, "1.10", nil)

	return c, teardown
}

func checkResponseSuccess(result *netapp.SingleResultBase, err error, t *testing.T) {
	if err != nil {
		t.Fatalf("Should not have gotten an error %s", err)
	}

	if !result.Passed() {
		t.Fatalf("Got the failure response, expected success, reason: %s", result.Reason)
	}
}

func checkAsyncResponseSuccess(result *netapp.AsyncResultBase, err error, t *testing.T) {
	if err != nil {
		t.Fatalf("Async Response should not have gotten an error %s", err)
	}

	if !result.Passed() {
		t.Fatalf("Async Response expected success got failure, reason: %s", result.ErrorMessage)
	}
}

func checkResponseFailure(result *netapp.SingleResultBase, err error, t *testing.T) {
	if err != nil {
		t.Fatalf("Should not have gotten an error %s", err)
	}

	if result.Passed() {
		t.Fatal("Got the successful response, expecting failure")
	}
}

func testFailureResult(errorNo int, reason string, result *netapp.SingleResultBase, t *testing.T) {

	if result.ErrorNo != errorNo {
		t.Errorf("%s got = %+v, want %+v", t.Name(), result.ErrorNo, errorNo)
	}

	if result.Reason != reason {
		t.Errorf("%s got = %+v, want %+v", t.Name(), result.Reason, reason)
	}
}

// debugTableItems is used to get reflected vaules of 2 items so its easier to tell why reflect.DeepEqual() fails
func debugItems(v1 interface{}, v2 interface{}) {
	val1 := reflect.ValueOf(v1)
	val2 := reflect.ValueOf(v2)
	fmt.Printf("v1: %v, v2: %v", val1, val2)
}