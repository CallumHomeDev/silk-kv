package datastore

import (
    "strconv"
    "testing"
)

func BenchmarkStore_SetGet_Sequential(b *testing.B) {
    s := New(0)
    for i := 0; i < b.N; i++ {
        key := "key" + strconv.Itoa(i)
        s.Set(key, []byte("value"))
        s.Get(key)
    }
}

func BenchmarkStore_SetGet_Parallel(b *testing.B) {
    s := New(0)
    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            key := "key" + strconv.Itoa(i)
            s.Set(key, []byte("value"))
            s.Get(key)
            i++
        }
    })
}