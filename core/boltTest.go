package core

import (
	"github.com/boltdb/bolt"
	"fmt"
)

//连接数据库
func TestBolt() {
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("bucketName"))
		if bucket == nil { //如果为空的话就创建这个bucket
			bucket, err = tx.CreateBucket([]byte("bucketName"))
			if err != nil {
				return err
			}
		}
		//写入数据
		bucket.Put([]byte("name1"), []byte("张三"))
		bucket.Put([]byte("name2"), []byte("李四"))
		//读取数据
		name1 := bucket.Get([]byte("name1"))
		name2 := bucket.Get([]byte("name2"))
		fmt.Printf("name1 is %s,name2 is %s\n", name1, name2)
		return nil
	})
}
