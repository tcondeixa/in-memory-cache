package cache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type CacheTestNestedValue struct {
	Field61 string
	Field62 int
	Field63 bool
}

type CacheTestValue struct {
	Field1 string
	Field2 int
	Field3 bool
	Field4 map[string]string
	Field5 []string
	Field6 map[string]CacheTestNestedValue
}

func (c *CacheTestValue) Clone() *CacheTestValue {
	cache := CacheTestValue{
		Field1: c.Field1,
		Field2: c.Field2,
		Field3: c.Field3,
		Field4: map[string]string{},
		Field5: []string{},
		Field6: map[string]CacheTestNestedValue{},
	}

	for k, v := range c.Field4 {
		cache.Field4[k] = v
	}

	cache.Field5 = append(cache.Field5, c.Field5...)

	for k, v := range c.Field6 {
		cache.Field6[k] = v
	}

	return &cache

}

var setSamples map[string]CacheTestValue

func init() {

	setSamples = map[string]CacheTestValue{
		"0": {
			Field1: "a",
			Field2: 1,
			Field3: true,
			Field4: map[string]string{
				"a": "a",
				"b": "b",
			},
			Field5: []string{
				"a",
				"b",
				"c",
			},
			Field6: map[string]CacheTestNestedValue{
				"1": {
					Field61: "a",
					Field62: 1,
					Field63: true,
				},
				"2": {
					Field61: "b",
					Field62: 2,
					Field63: false,
				},
			},
		},
		"1": {
			Field1: "a",
			Field2: 1,
			Field3: true,
			Field4: map[string]string{
				"a": "a",
				"b": "b",
			},
			Field5: []string{
				"a",
				"b",
				"c",
			},
			Field6: map[string]CacheTestNestedValue{
				"1": {
					Field61: "a",
					Field62: 1,
					Field63: true,
				},
				"2": {
					Field61: "b",
					Field62: 2,
					Field63: false,
				},
			},
		},
	}
}

// Benchmark Test
func BenchmarkInsert(b *testing.B) {
	cache := NewCache[CacheTestValue](time.Duration(3) * time.Second)

	for n := 0; n < b.N; n++ {
		err := cache.Insert(fmt.Sprintf("%d", n), CacheTestValue{})
		if err != nil {
			log.Fatal("Insert during benchmark")
		}
	}
}

func BenchmarkInsertWithClone(b *testing.B) {
	cache := NewCache[CacheTestValue](time.Duration(3) * time.Second)

	for n := 0; n < b.N; n++ {
		sample := setSamples["0"]
		err := cache.Insert(fmt.Sprintf("%d", n), *sample.Clone())
		if err != nil {
			log.Fatal("Insert during benchmark")
		}
	}
}

func BenchmarkUpsert(b *testing.B) {
	cache := NewCache[CacheTestValue](time.Duration(3) * time.Second)

	for n := 0; n < b.N; n++ {
		cache.Upsert(fmt.Sprintf("%d", n), CacheTestValue{})
	}
}

func BenchmarkBulkUpsert(b *testing.B) {
	cache := NewCache[CacheTestValue](time.Duration(3) * time.Second)
	all := map[string]CacheTestValue{}

	for n := 0; n < 100; n++ {
		all[fmt.Sprintf("%d", n)] = CacheTestValue{}
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		cache.BulkUpsert(all)
	}
}

func BenchmarkBulkInsert(b *testing.B) {
	cache := NewCache[CacheTestValue](time.Duration(3) * time.Second)
	all := map[int]map[string]CacheTestValue{}

	for n := 0; n < b.N; n++ {
		all[n] = map[string]CacheTestValue{}
		for i := 0; i < 100; i++ {
			all[n][fmt.Sprintf("%d", (n*100)+i)] = CacheTestValue{}
		}
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, err := cache.BulkInsert(all[n])
		if err != nil {
			log.Fatalf("BulkInsert during benchmark: %v", err)
		}
	}
}

// Unit Tests
func TestLoad(t *testing.T) {

	cache := NewCache[CacheTestValue](time.Duration(3) * time.Second)

	for _, tc := range []struct {
		testName      string
		filePath      string
		fileFormat    FileFormat
		expectedError error
	}{
		{
			testName:      "load file json",
			filePath:      "test_data/cache_dump.json",
			fileFormat:    Json,
			expectedError: nil,
		},
		{
			testName:      "load file yaml",
			filePath:      "test_data/cache_dump.yaml",
			fileFormat:    Yaml,
			expectedError: nil,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			err := cache.Load(tc.filePath, tc.fileFormat)
			require.Equal(t, tc.expectedError, err)
		})
	}
}

func TestSave(t *testing.T) {

	cache := NewCache[CacheTestValue](time.Duration(3) * time.Second)

	// Fill the Cache first
	for k, v := range setSamples {
		err := cache.Insert(k, v)
		if err != nil {
			log.Fatal("Insert during TestSave")
		}
	}

	for _, tc := range []struct {
		testName      string
		filePath      string
		fileFormat    FileFormat
		expectedError error
	}{
		{
			testName:      "save file json",
			filePath:      "/tmp/cache_dump.json",
			fileFormat:    Json,
			expectedError: nil,
		},
		{
			testName:      "save file yaml",
			filePath:      "/tmp/cache_dump.yaml",
			fileFormat:    Yaml,
			expectedError: nil,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			err := cache.Save(tc.filePath, tc.fileFormat)
			require.Equal(t, tc.expectedError, err)
		})
	}
}

func TestInsert(t *testing.T) {

	cache := NewCache[CacheTestValue](time.Duration(3) * time.Second)

	for _, tc := range []struct {
		testName      string
		sampleKey     string
		sampleValue   CacheTestValue
		expectedError error
	}{
		{
			testName:      "insert 0",
			sampleKey:     "0",
			sampleValue:   setSamples["0"],
			expectedError: nil,
		},
		{
			testName:      "insert 1",
			sampleKey:     "1",
			sampleValue:   setSamples["1"],
			expectedError: nil,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			err := cache.Insert(tc.sampleKey, tc.sampleValue)
			require.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpsert(t *testing.T) {

	cache := NewCache[CacheTestValue](time.Duration(3) * time.Second)

	for _, tc := range []struct {
		testName      string
		sampleKey     string
		sampleValue   CacheTestValue
		expectedError error
	}{
		{
			testName:      "upsert new 0",
			sampleKey:     "0",
			sampleValue:   setSamples["0"],
			expectedError: nil,
		},
		{
			testName:      "upsert new 1",
			sampleKey:     "1",
			sampleValue:   setSamples["1"],
			expectedError: nil,
		},
		{
			testName:      "upsert existing key 0",
			sampleKey:     "0",
			sampleValue:   setSamples["1"],
			expectedError: nil,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			cache.Upsert(tc.sampleKey, tc.sampleValue)
			value, err := cache.Get(tc.sampleKey)
			require.Equal(t, err, tc.expectedError)
			require.True(t, reflect.DeepEqual(*value, tc.sampleValue))
		})
	}
}

func TestGet(t *testing.T) {

	cache := NewCache[CacheTestValue](time.Duration(1) * time.Second)

	// Fill the Cache first
	for k, v := range setSamples {
		err := cache.Insert(k, v)
		if err != nil {
			log.Fatal("Insert during TestGet")
		}
	}

	for _, tc := range []struct {
		testName      string
		sampleKey     string
		sleep         time.Duration
		expectedValue CacheTestValue
		expectedError error
	}{
		{
			testName:      "get 0",
			sampleKey:     "0",
			sleep:         time.Second * time.Duration(0),
			expectedValue: setSamples["0"],
			expectedError: nil,
		},
		{
			testName:      "get 1",
			sampleKey:     "1",
			sleep:         time.Second * time.Duration(0),
			expectedValue: setSamples["1"],
			expectedError: nil,
		},
		{
			testName:      "get non-existent key 2",
			sampleKey:     "2",
			sleep:         time.Second * time.Duration(0),
			expectedValue: CacheTestValue{},
			expectedError: &InvalidCacheKeyError{},
		},
		{
			testName:      "get expired key 1",
			sampleKey:     "1",
			sleep:         time.Second * time.Duration(2),
			expectedValue: CacheTestValue{},
			expectedError: &OutdatedCacheEntryError{},
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			time.Sleep(tc.sleep)
			value, err := cache.Get(tc.sampleKey)
			require.ErrorIs(t, tc.expectedError, err)
			if err == nil {
				require.True(t, reflect.DeepEqual(*value, tc.expectedValue))
			}

		})
	}
}

func TestSize(t *testing.T) {

	cache := NewCache[CacheTestValue](time.Duration(3) * time.Second)

	for _, tc := range []struct {
		testName      string
		sampleKey     string
		sampleValue   CacheTestValue
		expectedSize  int
		expectedError error
	}{
		{
			testName:      "insert 0",
			sampleKey:     "0",
			sampleValue:   setSamples["0"],
			expectedSize:  1,
			expectedError: nil,
		},
		{
			testName:      "insert 1",
			sampleKey:     "1",
			sampleValue:   setSamples["1"],
			expectedSize:  2,
			expectedError: nil,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			err := cache.Insert(tc.sampleKey, tc.sampleValue)
			require.ErrorIs(t, tc.expectedError, err)
			require.Equal(t, tc.expectedSize, cache.Size())
		})
	}
}

func TestFlush(t *testing.T) {

	cache := NewCache[CacheTestValue](time.Duration(3) * time.Second)

	// Fill the Cache first
	for k, v := range setSamples {
		err := cache.Insert(k, v)
		if err != nil {
			log.Fatal("Insert during TestExpire")
		}
	}

	for _, tc := range []struct {
		testName     string
		expectedSize int
	}{
		{
			testName:     "flush cache non empty",
			expectedSize: 0,
		},
		{
			testName:     "flush again when empty",
			expectedSize: 0,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			cache.Flush()
			require.Equal(t, tc.expectedSize, cache.Size())
		})
	}
}

func TestDelete(t *testing.T) {

	cache := NewCache[CacheTestValue](time.Duration(3) * time.Second)

	// Fill the Cache first
	for k, v := range setSamples {
		err := cache.Insert(k, v)
		if err != nil {
			log.Fatal("Insert during TestExpire")
		}
	}

	for _, tc := range []struct {
		testName string
		key      string
	}{
		{
			testName: "delete first",
			key:      "1",
		},
		{
			testName: "delete second",
			key:      "2",
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			cache.Delete(tc.key)
			_, err := cache.Get(tc.key)
			require.ErrorIs(t, &InvalidCacheKeyError{}, err)
		})
	}
}

func TestExpire(t *testing.T) {

	cache := NewCache[CacheTestValue](time.Duration(1) * time.Second)

	// Fill the Cache first
	for k, v := range setSamples {
		err := cache.Insert(k, v)
		if err != nil {
			log.Fatal("Insert during TestExpire")
		}
	}

	for _, tc := range []struct {
		testName     string
		key          string
		sleep        time.Duration
		expectedSize int
	}{
		{
			testName:     "delete first",
			key:          "1",
			sleep:        time.Second * time.Duration(0),
			expectedSize: 2,
		},
		{
			testName:     "delete second",
			key:          "2",
			sleep:        time.Second * time.Duration(2),
			expectedSize: 0,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			time.Sleep(tc.sleep)
			cache.Expire()
			require.Equal(t, tc.expectedSize, cache.Size())
		})
	}
}
