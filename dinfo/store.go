package main

import "github.com/peterbourgon/diskv"

var d *diskv.Diskv

func init()  {
	// Simplest transform function: put all the data files into the base dir.
	flatTransform := func(s string) []string { return []string{} }

	// Initialize a new diskv store, rooted at "my-data-dir", with a 1MB cache.
	d = diskv.New(diskv.Options{
		BasePath:     "my-data-dir",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
}

func set(k, v string) error {
	return d.Write(k, []byte(v))
}

func get(k string) ([]byte, error) {
	val, err := d.Read(k)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func keys() []string {
	cancel := make(chan struct{})
	keysC := d.Keys(cancel)

	var keyList []string
	for key := range keysC {
		keyList = append(keyList, key)
	}

	return keyList
}