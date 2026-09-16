// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"codeberg.org/tiny-frameworks/nexutils/cache"
)

var Caches map[string]*cache.LRUCache

func loadCaches() {

	Caches = make(map[string]*cache.LRUCache)
	Caches["seller"] = cache.New(
		SystemParams.Caches.Seller.Capacity,
		SystemParams.Caches.Seller.TTL,
		SystemParams.Caches.Seller.Cleanup)

	Caches["country"] = cache.New(
		SystemParams.Caches.Country.Capacity,
		SystemParams.Caches.Country.TTL,
		SystemParams.Caches.Country.Cleanup)

	Caches["currency"] = cache.New(
		SystemParams.Caches.Currency.Capacity,
		SystemParams.Caches.Currency.TTL,
		SystemParams.Caches.Currency.Cleanup)

	Caches["unit"] = cache.New(
		SystemParams.Caches.Unit.Capacity,
		SystemParams.Caches.Unit.TTL,
		SystemParams.Caches.Unit.Cleanup)

	Caches["payment"] = cache.New(
		SystemParams.Caches.Payment.Capacity,
		SystemParams.Caches.Payment.TTL,
		SystemParams.Caches.Payment.Cleanup)
}

// StopCaches cleanly terminates the background cleanup goroutines of all caches.
func StopCaches() {
	for _, c := range Caches {
		if c != nil {
			c.StopCleanup()
		}
	}
}
