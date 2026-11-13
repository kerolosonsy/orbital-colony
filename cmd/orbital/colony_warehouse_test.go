package main

import (
	"context"
	"strings"
	"testing"
)

func TestWarehouseReferenceTransactionsAndPersistenceOffline(t *testing.T) {
	files := readColonyFixture(t, "s29-transfer")
	files["warehouse_test.go"] = `package main
import("bytes";"encoding/json";"fmt";"net/http";"net/http/httptest";"os";"path/filepath";"sync";"testing")
func call(t *testing.T,mux *http.ServeMux,method,path,body string)*httptest.ResponseRecorder{
 t.Helper();w:=httptest.NewRecorder();mux.ServeHTTP(w,httptest.NewRequest(method,path,bytes.NewBufferString(body)));return w
}
func TestWarehouse(t *testing.T){
 t.Setenv("ORBITAL_DB_PATH",filepath.Join(t.TempDir(),"warehouse.sqlite"))
 mux:=http.NewServeMux();setup(mux)
 batch:="{\"id\":\"one\",\"goods\":{\"moss\":12.5,\"lichen\":3}}"
 for range 2{if w:=call(t,mux,"POST","/inventory",batch);w.Code!=204{t.Fatal(w.Code,w.Body)}}
 for range 3{mux=http.NewServeMux();setup(mux);if w:=call(t,mux,"POST","/inventory",batch);w.Code!=204{t.Fatal(w.Code,w.Body)};w:=call(t,mux,"GET","/inventory","");var inv map[string]float64;if json.Unmarshal(w.Body.Bytes(),&inv)!=nil || inv["moss"]!=12.5 || inv["lichen"]!=3{t.Fatal(w.Code,w.Body)}}
 var workers sync.WaitGroup
 for index:=range 40{workers.Add(1);go func(index int){defer workers.Done();id:=fmt.Sprintf("%d",index/4);w:=call(t,mux,"POST","/transfer",fmt.Sprintf("{\"id\":%q,\"from\":\"a\",\"to\":\"b\",\"amount\":7}",id));if w.Code!=200{t.Error(w.Code,w.Body)}}(index)}
 workers.Wait()
 w:=call(t,mux,"GET","/stores","");var stores map[string]int;json.Unmarshal(w.Body.Bytes(),&stores);if stores["a"]!=30 || stores["b"]!=170{t.Fatal(stores)}
 for index:=range 20{workers.Add(1);go func(index int){defer workers.Done();w:=call(t,mux,"POST","/transfer",fmt.Sprintf("{\"id\":\"overdraw-%d\",\"from\":\"a\",\"to\":\"b\",\"amount\":31}",index));if w.Code!=409{t.Error(w.Code,w.Body)}}(index)}
 workers.Wait();w=call(t,mux,"GET","/stores","");json.Unmarshal(w.Body.Bytes(),&stores);if stores["a"]!=30 || stores["b"]!=170{t.Fatal(stores)}
 header,err:=os.ReadFile(os.Getenv("ORBITAL_DB_PATH"));if err!=nil || string(header[:16])!="SQLite format 3\x00"{t.Fatal(err)}
}
`
	passed, output, err := testRunner(t).goTest(context.Background(), colonyControlLevel(), files, "-race", "./...")
	if err != nil || !passed {
		t.Fatalf("offline warehouse reference: %v\n%s", err, output)
	}
	if !strings.Contains(output, "ok") {
		t.Fatal(output)
	}
}
