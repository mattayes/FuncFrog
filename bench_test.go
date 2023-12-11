package funcfrog

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"

	"github.com/koss-null/funcfrog/pkg/ff"
)

type User struct {
	ID string
	// other stuff here
}

func GetUserID(u *User) string {
	return u.ID
}

func BenchmarkSerial(b *testing.B) {
	// n == number of users
	for _, n := range []int{
		1,
		100,
		10_000,
		1_000_000,
		100_000_000,
	} {
		b.Run(fmt.Sprintf("std-%d-users", n), func(b *testing.B) {
			users := makeUsers(n)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ids := make([]string, n)
				for i, user := range users {
					ids[i] = user.ID
				}
			}
		})
		b.Run(fmt.Sprintf("funcfrog-%d-users", n), func(b *testing.B) {
			users := makeUsers(n)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = ff.Map(users, GetUserID).Do()
			}
		})
	}
}

func BenchmarkConcurrent(b *testing.B) {
	// n == number of users
	maxWorkers := runtime.GOMAXPROCS(0)
	for _, n := range []int{
		1,
		100,
		10_000,
		1_000_000,
		100_000_000,
	} {
		numWorkers := maxWorkers
		if n < maxWorkers {
			// No sense in having more workers than items
			numWorkers = n
		}

		b.Run(fmt.Sprintf("std-%d-users", n), func(b *testing.B) {
			users := makeUsers(n)
			b.ResetTimer()

			chunkSize := n / numWorkers
			for i := 0; i < b.N; i++ {
				ids := make([]string, n)
				var wg sync.WaitGroup
				wg.Add(numWorkers)
				for i := 0; i < numWorkers; i++ {
					go func(i int) {
						start := i * n / chunkSize
						end := (i + 1) * n / chunkSize
						if end > n {
							// out-of-bounds
							end = n
						}
						for i := start; i < end; i++ {
							ids[i] = users[i].ID
						}
						wg.Done()
					}(i)
				}
				wg.Wait()
			}
		})
		b.Run(fmt.Sprintf("funcfrog-%d-users", n), func(b *testing.B) {
			users := makeUsers(n)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = ff.Map(users, GetUserID).Parallel(uint16(numWorkers)).Do()
			}
		})
	}
}

func makeUsers(n int) []*User {
	users := make([]*User, n)
	for i := 0; i < n; i++ {
		users[i] = &User{ID: strconv.Itoa(i)}
	}
	return users
}
