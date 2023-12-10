package funcfrog^$''

import (
	"strconv"
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

func BenchmarkAllTheThings(b *testing.B) {
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

func makeUsers(n int) []*User {
	users := make([]*User, n)
	for i := 0; i < n; i++ {
		users[i] = &User{ID: strconv.Itoa(i)}
	}
	return users
}
