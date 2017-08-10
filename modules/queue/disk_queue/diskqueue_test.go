/*
[nsq]: https://github.com/nsqio/nsq
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
 [nsq]: https://github.com/nsqio/nsq
*/
package queue

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"testing"
	"time"
)

func TestDiskQueue(t *testing.T) {
	dqName := "test_disk_queue" + strconv.Itoa(int(time.Now().Unix()))
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("nsq-test-%d", time.Now().UnixNano()))
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)
	dq := NewDiskQueue(dqName, tmpDir, 1024, 4, 1<<10, 2500, 2*time.Second)
	defer dq.Close()
	assert.NotEqual(t, nil, dq)
	assert.Equal(t, dq.Depth(), int64(0))

	msg := []byte("test")
	err = dq.Put(msg)
	assert.Equal(t, err, nil)
	assert.Equal(t, dq.Depth(), int64(1))

	msgOut := <-dq.ReadChan()
	assert.Equal(t, msgOut, msg)
}

func BenchmarkDiskQueuePut16(b *testing.B) {
	benchmarkDiskQueuePut(16, b)
}
func BenchmarkDiskQueuePut64(b *testing.B) {
	benchmarkDiskQueuePut(64, b)
}
func BenchmarkDiskQueuePut256(b *testing.B) {
	benchmarkDiskQueuePut(256, b)
}
func BenchmarkDiskQueuePut1024(b *testing.B) {
	benchmarkDiskQueuePut(1024, b)
}
func BenchmarkDiskQueuePut4096(b *testing.B) {
	benchmarkDiskQueuePut(4096, b)
}
func BenchmarkDiskQueuePut16384(b *testing.B) {
	benchmarkDiskQueuePut(16384, b)
}
func BenchmarkDiskQueuePut65536(b *testing.B) {
	benchmarkDiskQueuePut(65536, b)
}
func BenchmarkDiskQueuePut262144(b *testing.B) {
	benchmarkDiskQueuePut(262144, b)
}
func BenchmarkDiskQueuePut1048576(b *testing.B) {
	benchmarkDiskQueuePut(1048576, b)
}
func benchmarkDiskQueuePut(size int64, b *testing.B) {
	b.StopTimer()
	dqName := "bench_disk_queue_put" + strconv.Itoa(b.N) + strconv.Itoa(int(time.Now().Unix()))
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("nsq-test-%d", time.Now().UnixNano()))
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)
	dq := NewDiskQueue(dqName, tmpDir, 1024768*100, 0, 1<<20, 2500, 2*time.Second)
	defer dq.Close()
	b.SetBytes(size)
	data := make([]byte, size)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		err := dq.Put(data)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkDiskWrite16(b *testing.B) {
	benchmarkDiskWrite(16, b)
}
func BenchmarkDiskWrite64(b *testing.B) {
	benchmarkDiskWrite(64, b)
}
func BenchmarkDiskWrite256(b *testing.B) {
	benchmarkDiskWrite(256, b)
}
func BenchmarkDiskWrite1024(b *testing.B) {
	benchmarkDiskWrite(1024, b)
}
func BenchmarkDiskWrite4096(b *testing.B) {
	benchmarkDiskWrite(4096, b)
}
func BenchmarkDiskWrite16384(b *testing.B) {
	benchmarkDiskWrite(16384, b)
}
func BenchmarkDiskWrite65536(b *testing.B) {
	benchmarkDiskWrite(65536, b)
}
func BenchmarkDiskWrite262144(b *testing.B) {
	benchmarkDiskWrite(262144, b)
}
func BenchmarkDiskWrite1048576(b *testing.B) {
	benchmarkDiskWrite(1048576, b)
}
func benchmarkDiskWrite(size int64, b *testing.B) {
	b.StopTimer()
	fileName := "bench_disk_queue_put" + strconv.Itoa(b.N) + strconv.Itoa(int(time.Now().Unix()))
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("nsq-test-%d", time.Now().UnixNano()))
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)
	f, _ := os.OpenFile(path.Join(tmpDir, fileName), os.O_RDWR|os.O_CREATE, 0600)
	b.SetBytes(size)
	data := make([]byte, size)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		f.Write(data)
	}
	f.Sync()
}

func BenchmarkDiskWriteBuffered16(b *testing.B) {
	benchmarkDiskWriteBuffered(16, b)
}
func BenchmarkDiskWriteBuffered64(b *testing.B) {
	benchmarkDiskWriteBuffered(64, b)
}
func BenchmarkDiskWriteBuffered256(b *testing.B) {
	benchmarkDiskWriteBuffered(256, b)
}
func BenchmarkDiskWriteBuffered1024(b *testing.B) {
	benchmarkDiskWriteBuffered(1024, b)
}
func BenchmarkDiskWriteBuffered4096(b *testing.B) {
	benchmarkDiskWriteBuffered(4096, b)
}
func BenchmarkDiskWriteBuffered16384(b *testing.B) {
	benchmarkDiskWriteBuffered(16384, b)
}
func BenchmarkDiskWriteBuffered65536(b *testing.B) {
	benchmarkDiskWriteBuffered(65536, b)
}
func BenchmarkDiskWriteBuffered262144(b *testing.B) {
	benchmarkDiskWriteBuffered(262144, b)
}
func BenchmarkDiskWriteBuffered1048576(b *testing.B) {
	benchmarkDiskWriteBuffered(1048576, b)
}
func benchmarkDiskWriteBuffered(size int64, b *testing.B) {
	b.StopTimer()
	fileName := "bench_disk_queue_put" + strconv.Itoa(b.N) + strconv.Itoa(int(time.Now().Unix()))
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("nsq-test-%d", time.Now().UnixNano()))
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)
	f, _ := os.OpenFile(path.Join(tmpDir, fileName), os.O_RDWR|os.O_CREATE, 0600)
	b.SetBytes(size)
	data := make([]byte, size)
	w := bufio.NewWriterSize(f, 1024*4)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		w.Write(data)
		if i%1024 == 0 {
			w.Flush()
		}
	}
	w.Flush()
	f.Sync()
}

// you might want to run this like
// $ go test -bench=DiskQueueGet -benchtime 0.1s
// too avoid doing too many iterations.
func BenchmarkDiskQueueGet16(b *testing.B) {
	benchmarkDiskQueueGet(16, b)
}
func BenchmarkDiskQueueGet64(b *testing.B) {
	benchmarkDiskQueueGet(64, b)
}
func BenchmarkDiskQueueGet256(b *testing.B) {
	benchmarkDiskQueueGet(256, b)
}
func BenchmarkDiskQueueGet1024(b *testing.B) {
	benchmarkDiskQueueGet(1024, b)
}
func BenchmarkDiskQueueGet4096(b *testing.B) {
	benchmarkDiskQueueGet(4096, b)
}
func BenchmarkDiskQueueGet16384(b *testing.B) {
	benchmarkDiskQueueGet(16384, b)
}
func BenchmarkDiskQueueGet65536(b *testing.B) {
	benchmarkDiskQueueGet(65536, b)
}
func BenchmarkDiskQueueGet262144(b *testing.B) {
	benchmarkDiskQueueGet(262144, b)
}
func BenchmarkDiskQueueGet1048576(b *testing.B) {
	benchmarkDiskQueueGet(1048576, b)
}

func benchmarkDiskQueueGet(size int64, b *testing.B) {
	b.StopTimer()
	dqName := "bench_disk_queue_get" + strconv.Itoa(b.N) + strconv.Itoa(int(time.Now().Unix()))
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("nsq-test-%d", time.Now().UnixNano()))
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)
	dq := NewDiskQueue(dqName, tmpDir, 1024768, 0, 1<<10, 2500, 2*time.Second)
	defer dq.Close()
	b.SetBytes(size)
	data := make([]byte, size)
	for i := 0; i < b.N; i++ {
		dq.Put(data)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		<-dq.ReadChan()
	}
}
